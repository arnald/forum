package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/arnald/forum/internal/domain/categories"
)

type categoryRepo struct {
	DB *sql.DB
}

func NewCategoryRepo(db *sql.DB) *categoryRepo {
	return &categoryRepo{DB: db}
}

func (r *categoryRepo) CreateCategory(ctx context.Context, category *categories.Category) error {
	query := `
	INSERT INTO categories (name, description, created_by)
	VALUES (?,?,?)`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		category.Name,
		category.Description,
		category.CreatedBy,
	)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "UNIQUE constraint failed: categories.name"):
			return fmt.Errorf("category with name %s already exists: %w", category.Name, ErrCategoryAlreadyExists)
		case strings.Contains(err.Error(), "FOREIGN KEY constraint failed: categories.created_by"):
			return fmt.Errorf("user with ID %s not found: %w", category.CreatedBy, ErrUserNotFound)
		default:
			return fmt.Errorf("failed to create category: %w", err)
		}
	}
	return nil
}
