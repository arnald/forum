package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	// Need to import sqlite driver.
	_ "github.com/mattn/go-sqlite3"

	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/pkg/path"
)

const (
	permissionUserRWE = 0o750
)

func InitializeDB(cfg config.ServerConfig) (*sql.DB, error) {
	// Ensure directory exists
	err := os.MkdirAll(filepath.Dir(cfg.Database.Path), permissionUserRWE)
	if err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	db, result, err := OpenDB(cfg)
	if err != nil {
		return result, err
	}

	if cfg.Database.MigrateOnStart {
		err := migrateDB(db)
		if err != nil {
			return nil, fmt.Errorf("migration failed: %w", err)
		}
	}

	if cfg.Database.SeedOnStart && cfg.Database.MigrateOnStart {
		err := seedDB(db, cfg.Environment)
		if err != nil {
			log.Printf("Seeding warning: %v", err)
		}
	}

	return db, nil
}

func OpenDB(cfg config.ServerConfig) (*sql.DB, *sql.DB, error) {
	db, err := sql.Open(cfg.Database.Driver, cfg.Database.Path+"?"+cfg.Database.Pragma)
	if err != nil {
		return nil, nil, err
	}

	if cfg.Database.Driver == "sqlite3" {
		db.SetMaxOpenConns(cfg.Database.OpenConn)
	}
	return db, nil, nil
}

func migrateDB(db *sql.DB) error {
	resolver := path.NewResolver()
	migrationFiles := []string{
		resolver.GetPath("db/migrations/schema.sql"),
		resolver.GetPath("db/migrations/indexes.sql"),
	}

	for _, file := range migrationFiles {
		err := execSQLFile(db, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func execSQLFile(db *sql.DB, path string) error {
	ctx := context.TODO()

	absPath := filepath.Clean(path)
	if !filepath.IsAbs(absPath) {
		absPath = filepath.Join(".", absPath)
	}
	absPath = filepath.Clean(absPath)

	content, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %w", err)
	}

	transaction, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Handle rollback errors safely
	defer func() {
		rollbackErr := transaction.Rollback()
		if rollbackErr != nil {
			if !errors.Is(rollbackErr, sql.ErrTxDone) {
				if err == nil {
					err = fmt.Errorf("rollback failed: %w", rollbackErr)
				} else {
					err = errors.Join(
						err,
						fmt.Errorf("rollback failed: %w", rollbackErr),
					)
				}
			}
		}
	}()

	statements := strings.SplitSeq(string(content), ":")
	for stmt := range statements {
		trimmed := strings.TrimSpace(stmt)
		if trimmed == "" {
			continue
		}
		_, err = transaction.ExecContext(ctx, trimmed)
		if err != nil {
			return fmt.Errorf("error executing '%s....': %w", trimmed, err)
		}
	}

	// Commit transaction
	err = transaction.Commit()
	if err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}
	return nil
}

func seedDB(db *sql.DB, env string) error {
	resolver := path.NewResolver()
	switch env {
	case "development":
		return execSQLFile(db, resolver.GetPath("db/seeds/dev_data.sql"))
	case "staging":
		return execSQLFile(db, resolver.GetPath("db/seeds/test.sql"))
	default:
		return nil // No seeding in production
	}
}
