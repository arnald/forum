package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/arnald/forum/internal/pkg/helpers"
	"github.com/arnald/forum/internal/pkg/path"
)

const (
	readTimeout  = 5
	writeTimeout = 10
	idleTimeout  = 15
)

var (
	ErrMissingServerHost    = errors.New("missing SERVER_HOST in config")
	ErrServerPortNotInteger = errors.New("invalid SERVER_PORT: must be integer")
)

type ServerConfig struct {
	Host         string
	Port         string
	Environment  string
	APIContext   string
	Database     DatabaseConfig
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Driver         string
	Path           string
	Pragma         string
	MigrateOnStart bool
	SeedOnStart    bool
	OpenConn       int
}

func LoadConfig() (*ServerConfig, error) {
	resolver := path.NewResolver()
	envFile, _ := os.ReadFile(resolver.GetPath(".env"))
	envMap := helpers.ParseEnv(string(envFile))

	cfg := &ServerConfig{
		Host:         helpers.GetEnv("SERVER_HOST", envMap, "localhost"),
		Port:         helpers.GetEnv("SERVER_PORT", envMap, "8080"),
		Environment:  helpers.GetEnv("SERVER_ENVIRONMENT", envMap, "development"),
		APIContext:   helpers.GetEnv("API_CONTEXT", envMap, "/api/v1"),
		ReadTimeout:  helpers.GetEnvDuration("SERVER_READ_TIMEOUT", envMap, readTimeout),
		WriteTimeout: helpers.GetEnvDuration("SERVER_WRITE_TIMEOUT", envMap, writeTimeout),
		IdleTimeout:  helpers.GetEnvDuration("SERVER_IDLE_TIMEOUT", envMap, idleTimeout),
		Database: DatabaseConfig{
			Driver:         helpers.GetEnv("DB_DRIVER", envMap, "sqlite3"),
			Path:           resolver.GetPath(helpers.GetEnv("DB_PATH", envMap, "data/forum.db")),
			MigrateOnStart: helpers.GetEnvBool("DB_MIGRATE_ON_START", envMap, true),
			SeedOnStart:    helpers.GetEnvBool("DB_SEED_ON_START", envMap, true),
			Pragma:         helpers.GetEnv("DB_PRAGMA", envMap, "_foreign_keys=on&_journal_mode=WAL"),
			OpenConn:       helpers.GetEnvInt("DB_OPEN_CONN", envMap, 1),
		},
	}

	if cfg.Host == "" {
		return nil, ErrMissingServerHost
	}
	_, err := strconv.Atoi(strings.TrimPrefix(cfg.Port, ":"))
	if err != nil {
		return nil, ErrServerPortNotInteger
	}

	return cfg, nil
}
