package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/arnald/forum/internal/domain/category"
)

func (r Repo) CreateCategory(ctx context.Context, categoryData *category.Category) error {
	query := `
	INSERT INTO categories (id, name, description, created_at)
	VALUES (?, ?, ?, ?)`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		categoryData.ID,
		categoryData.Name,
		categoryData.Description,
		categoryData.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}

	return nil
}

func (r Repo) GetCategoryByID(ctx context.Context, id string) (*category.Category, error) {
	query := `
	SELECT c.id, c.name, c.description, c.created_at,
	       COALESCE(COUNT(pc.post_id), 0) as post_count
	FROM categories c
	LEFT JOIN post_categories pc ON c.id = pc.category_id
	WHERE c.id = ?
	GROUP BY c.id, c.name, c.description, c.created_at`

	var categoryData category.Category
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&categoryData.ID,
		&categoryData.Name,
		&categoryData.Description,
		&categoryData.CreatedAt,
		&categoryData.PostCount,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("category with ID %s not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get category by ID: %w", err)
	}

	return &categoryData, nil
}

func (r Repo) GetCategoryByName(ctx context.Context, name string) (*category.Category, error) {
	query := `
	SELECT c.id, c.name, c.description, c.created_at,
	       COALESCE(COUNT(pc.post_id), 0) as post_count
	FROM categories c
	LEFT JOIN post_categories pc ON c.id = pc.category_id
	WHERE c.name = ?
	GROUP BY c.id, c.name, c.description, c.created_at`

	var categoryData category.Category
	err := r.DB.QueryRowContext(ctx, query, name).Scan(
		&categoryData.ID,
		&categoryData.Name,
		&categoryData.Description,
		&categoryData.CreatedAt,
		&categoryData.PostCount,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("category with name %s not found", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get category by name: %w", err)
	}

	return &categoryData, nil
}

func (r Repo) GetAllCategories(ctx context.Context) ([]*category.Category, error) {
	query := `
	SELECT c.id, c.name, c.description, c.created_at,
	       COALESCE(COUNT(pc.post_id), 0) as post_count
	FROM categories c
	LEFT JOIN post_categories pc ON c.id = pc.category_id
	GROUP BY c.id, c.name, c.description, c.created_at
	ORDER BY c.name ASC`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []*category.Category
	for rows.Next() {
		var categoryData category.Category
		err := rows.Scan(
			&categoryData.ID,
			&categoryData.Name,
			&categoryData.Description,
			&categoryData.CreatedAt,
			&categoryData.PostCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, &categoryData)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return categories, nil
}

func (r Repo) UpdateCategory(ctx context.Context, categoryData *category.Category) error {
	query := `
	UPDATE categories
	SET name = ?, description = ?
	WHERE id = ?`

	result, err := r.DB.ExecContext(
		ctx,
		query,
		categoryData.Name,
		categoryData.Description,
		categoryData.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("category with ID %s not found", categoryData.ID)
	}

	return nil
}

func (r Repo) DeleteCategory(ctx context.Context, id string) error {
	// First check if category has posts
	countQuery := `SELECT COUNT(*) FROM post_categories WHERE category_id = ?`
	var postCount int
	err := r.DB.QueryRowContext(ctx, countQuery, id).Scan(&postCount)
	if err != nil {
		return fmt.Errorf("failed to check category post count: %w", err)
	}

	if postCount > 0 {
		return fmt.Errorf("cannot delete category with %d posts", postCount)
	}

	query := `DELETE FROM categories WHERE id = ?`
	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("category with ID %s not found", id)
	}

	return nil
}

func (r Repo) GetCategoryWithPosts(ctx context.Context, categoryID string, limit int, offset int) (*category.CategoryWithPosts, error) {
	// Get category info
	categoryData, err := r.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	// Get posts in this category
	postsQuery := `
	SELECT p.id, p.title, p.user_id, u.username, p.created_at
	FROM posts p
	JOIN post_categories pc ON p.id = pc.post_id
	JOIN users u ON p.user_id = u.id
	WHERE pc.category_id = ?
	ORDER BY p.created_at DESC
	LIMIT ? OFFSET ?`

	rows, err := r.DB.QueryContext(ctx, postsQuery, categoryID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query category posts: %w", err)
	}
	defer rows.Close()

	var posts []category.CategoryPost
	for rows.Next() {
		var post category.CategoryPost
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.UserID,
			&post.Username,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category post: %w", err)
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return &category.CategoryWithPosts{
		Category: categoryData,
		Posts:    posts,
	}, nil
}

func (r Repo) GetPostsCount(ctx context.Context, categoryID string) (int, error) {
	query := `SELECT COUNT(*) FROM post_categories WHERE category_id = ?`
	var count int
	err := r.DB.QueryRowContext(ctx, query, categoryID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get posts count for category: %w", err)
	}
	return count, nil
}