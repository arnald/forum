package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/arnald/forum/internal/domain/post"
	"github.com/arnald/forum/internal/domain/user"
)

func (r Repo) CreatePost(ctx context.Context, postData *post.Post) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
	INSERT INTO posts (id, title, content, user_id, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?)`

	_, err = tx.ExecContext(
		ctx,
		query,
		postData.ID,
		postData.Title,
		postData.Content,
		postData.UserID,
		postData.CreatedAt,
		postData.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert post: %w", err)
	}

	for _, category := range postData.Categories {
		_, err = tx.ExecContext(
			ctx,
			`INSERT OR IGNORE INTO post_categories (post_id, category_id) VALUES (?, ?)`,
			postData.ID,
			category,
		)
		if err != nil {
			return fmt.Errorf("failed to associate post with category: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r Repo) GetPostByID(ctx context.Context, id string) (*post.Post, error) {
	query := `
	SELECT p.id, p.title, p.content, p.user_id, u.username, p.created_at, p.updated_at
	FROM posts p
	JOIN users u ON p.user_id = u.id
	WHERE p.id = ?`

	var postData post.Post
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&postData.ID,
		&postData.Title,
		&postData.Content,
		&postData.UserID,
		&postData.Username,
		&postData.CreatedAt,
		&postData.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("post with ID %s not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get post by ID: %w", err)
	}

	categories, err := r.getPostCategories(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post categories: %w", err)
	}
	postData.Categories = categories

	return &postData, nil
}

func (r Repo) GetAllPosts(ctx context.Context) ([]*post.Post, error) {
	query := `
	SELECT p.id, p.title, p.content, p.user_id, u.username, p.created_at, p.updated_at
	FROM posts p
	JOIN users u ON p.user_id = u.id
	ORDER BY p.created_at DESC`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	var posts []*post.Post
	for rows.Next() {
		var postData post.Post
		err := rows.Scan(
			&postData.ID,
			&postData.Title,
			&postData.Content,
			&postData.UserID,
			&postData.Username,
			&postData.CreatedAt,
			&postData.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}

		categories, err := r.getPostCategories(ctx, postData.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get post categories: %w", err)
		}
		postData.Categories = categories

		posts = append(posts, &postData)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return posts, nil
}

func (r Repo) GetPostsByUserID(ctx context.Context, userID string) ([]*post.Post, error) {
	query := `
	SELECT p.id, p.title, p.content, p.user_id, u.username, p.created_at, p.updated_at
	FROM posts p
	JOIN users u ON p.user_id = u.id
	WHERE p.user_id = ?
	ORDER BY p.created_at DESC`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts by user ID: %w", err)
	}
	defer rows.Close()

	var posts []*post.Post
	for rows.Next() {
		var postData post.Post
		err := rows.Scan(
			&postData.ID,
			&postData.Title,
			&postData.Content,
			&postData.UserID,
			&postData.Username,
			&postData.CreatedAt,
			&postData.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}

		categories, err := r.getPostCategories(ctx, postData.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get post categories: %w", err)
		}
		postData.Categories = categories

		posts = append(posts, &postData)
	}

	return posts, nil
}

func (r Repo) GetPostsByCategory(ctx context.Context, category string) ([]*post.Post, error) {
	query := `
	SELECT p.id, p.title, p.content, p.user_id, u.username, p.created_at, p.updated_at
	FROM posts p
	JOIN users u ON p.user_id = u.id
	JOIN post_categories pc ON p.id = pc.post_id
	JOIN categories c ON pc.category_id = c.id
	WHERE c.name = ?
	ORDER BY p.created_at DESC`

	rows, err := r.DB.QueryContext(ctx, query, category)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts by category: %w", err)
	}
	defer rows.Close()

	var posts []*post.Post
	for rows.Next() {
		var postData post.Post
		err := rows.Scan(
			&postData.ID,
			&postData.Title,
			&postData.Content,
			&postData.UserID,
			&postData.Username,
			&postData.CreatedAt,
			&postData.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}

		categories, err := r.getPostCategories(ctx, postData.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get post categories: %w", err)
		}
		postData.Categories = categories

		posts = append(posts, &postData)
	}

	return posts, nil
}

func (r Repo) GetPostsByCategoryID(ctx context.Context, categoryID string) ([]*post.Post, error) {
	query := `
	SELECT p.id, p.title, p.content, p.user_id, u.username, p.created_at, p.updated_at
	FROM posts p
	JOIN users u ON p.user_id = u.id
	JOIN post_categories pc ON p.id = pc.post_id
	WHERE pc.category_id = ?
	ORDER BY p.created_at DESC`

	rows, err := r.DB.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts by category ID: %w", err)
	}
	defer rows.Close()

	var posts []*post.Post
	for rows.Next() {
		var postData post.Post
		err := rows.Scan(
			&postData.ID,
			&postData.Title,
			&postData.Content,
			&postData.UserID,
			&postData.Username,
			&postData.CreatedAt,
			&postData.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}

		categories, err := r.getPostCategories(ctx, postData.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get post categories: %w", err)
		}
		postData.Categories = categories

		posts = append(posts, &postData)
	}

	return posts, nil
}

func (r Repo) UpdatePost(ctx context.Context, postData *post.Post) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
	UPDATE posts
	SET title = ?, content = ?, updated_at = ?
	WHERE id = ?`

	result, err := tx.ExecContext(
		ctx,
		query,
		postData.Title,
		postData.Content,
		time.Now(),
		postData.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("post with ID %s not found", postData.ID)
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM post_categories WHERE post_id = ?`, postData.ID)
	if err != nil {
		return fmt.Errorf("failed to delete existing post categories: %w", err)
	}

	for _, category := range postData.Categories {
		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`,
			postData.ID,
			category,
		)
		if err != nil {
			return fmt.Errorf("failed to associate post with category: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r Repo) DeletePost(ctx context.Context, id string) error {
	query := `DELETE FROM posts WHERE id = ?`
	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("post with ID %s not found", id)
	}

	return nil
}

func (r Repo) GetCategories(ctx context.Context) ([]*post.Category, error) {
	query := `SELECT id, name FROM categories ORDER BY name`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []*post.Category
	for rows.Next() {
		var category post.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, &category)
	}

	return categories, nil
}

func (r Repo) getPostCategories(ctx context.Context, postID string) ([]string, error) {
	query := `
	SELECT c.name
	FROM categories c
	JOIN post_categories pc ON c.id = pc.category_id
	WHERE pc.post_id = ?`

	rows, err := r.DB.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to query post categories: %w", err)
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r Repo) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	query := `
	SELECT id, username, email, password_hash, created_at, avatar_url
	FROM users
	WHERE id = ?`

	var userData user.User
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&userData.ID,
		&userData.Username,
		&userData.Email,
		&userData.Password,
		&userData.CreatedAt,
		&userData.AvatarURL,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user with ID %s not found: %w", id, ErrUserNotFound)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &userData, nil
}