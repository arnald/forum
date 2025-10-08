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

func (r *categoryRepo) GetAllCategories(ctx context.Context) ([]*categories.Category, error) {
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

	var categoriesList []*categories.Category
	for rows.Next() {
		var category categories.Category
		err := rows.Scan(
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
	return categoriesList, nil
}
func (r *categoryRepo) GetCategoryByID(ctx context.Context, id int) (*categories.Category, error) {
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

	var category categories.Category
	err = stmt.QueryRowContext(ctx, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedBy,
		&category.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category with ID %d not found: %w", id, ErrCategoryNotFound)
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return &category, nil
}

func (r *categoryRepo) DeleteCategory(ctx context.Context, id int) error {
	query := `
	DELETE FROM categories
	WHERE id = ?
	`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id)
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

func (r *categoryRepo) UpdateCategory(ctx context.Context, category *categories.Category) error {
	query := `
	UPDATE categories
	SET name = ?, description = ?, updated_at = CURRENT_TIMESTAMP
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
