package categories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/arnald/forum/internal/domain/category"
)

type Repo struct {
	DB *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{DB: db}
}

func (r *Repo) CreateCategory(ctx context.Context, category *category.Category) error {
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

func (r *Repo) GetAllCategories(ctx context.Context) ([]*category.Category, error) {
	query := `
	SELECT id, name, description, created_by, created_at
	FROM categories
	`
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	categoriesList := make([]*category.Category, 0)
	for rows.Next() {
		var category category.Category
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.CreatedBy,
			&category.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		categoriesList = append(categoriesList, &category)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return categoriesList, nil
}

func (r *Repo) GetCategoryByID(ctx context.Context, id int) (*category.Category, error) {
	query := `
	SELECT id, name, description, created_by, created_at
	FROM categories
	WHERE id = ?
	`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	var category category.Category
	err = stmt.QueryRowContext(ctx, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedBy,
		&category.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("category with ID %d not found: %w", id, ErrCategoryNotFound)
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return &category, nil
}

func (r *Repo) DeleteCategory(ctx context.Context, id int, userID string) error {
	query := `
	DELETE FROM categories
	WHERE id = ? AND created_by = ?
	`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("retrieving rows affected failed: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("category with ID %d not found: %w", id, ErrCategoryNotFound)
	}
	return nil
}

func (r *Repo) UpdateCategory(ctx context.Context, category *category.Category) error {
	query := `
	UPDATE categories
	SET name = ?, description = ?
	WHERE id = ?
	`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		category.Name,
		category.Description,
		category.ID,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: categories.name") {
			return fmt.Errorf("category with name %s already exists: %w", category.Name, ErrCategoryAlreadyExists)
		}
		return fmt.Errorf("exec failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("retrieving rows affected failed: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("category with ID %d not found: %w", category.ID, ErrCategoryNotFound)
	}
	return nil
}
