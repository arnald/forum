package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/arnald/forum/internal/domain/comment"
)

func (r Repo) CreateComment(ctx context.Context, commentData *comment.Comment) error {
	query := `
	INSERT INTO comments (id, content, post_id, user_id, parent_id, level, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		commentData.ID,
		commentData.Content,
		commentData.PostID,
		commentData.UserID,
		commentData.ParentID,
		commentData.Level,
		commentData.CreatedAt,
		commentData.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	return nil
}

func (r Repo) GetCommentByID(ctx context.Context, id string) (*comment.Comment, error) {
	query := `
	SELECT c.id, c.content, c.post_id, c.user_id, u.username, c.parent_id, c.level, c.created_at, c.updated_at
	FROM comments c
	JOIN users u ON c.user_id = u.id
	WHERE c.id = ?`

	var commentData comment.Comment
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&commentData.ID,
		&commentData.Content,
		&commentData.PostID,
		&commentData.UserID,
		&commentData.Username,
		&commentData.ParentID,
		&commentData.Level,
		&commentData.CreatedAt,
		&commentData.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("comment with ID %s not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get comment by ID: %w", err)
	}

	return &commentData, nil
}

func (r Repo) GetCommentsByPostID(ctx context.Context, postID string) ([]*comment.Comment, error) {
	query := `
	SELECT c.id, c.content, c.post_id, c.user_id, u.username, c.parent_id, c.level, c.created_at, c.updated_at
	FROM comments c
	JOIN users u ON c.user_id = u.id
	WHERE c.post_id = ?
	ORDER BY c.created_at ASC`

	rows, err := r.DB.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments by post ID: %w", err)
	}
	defer rows.Close()

	var comments []*comment.Comment
	for rows.Next() {
		var commentData comment.Comment
		err := rows.Scan(
			&commentData.ID,
			&commentData.Content,
			&commentData.PostID,
			&commentData.UserID,
			&commentData.Username,
			&commentData.ParentID,
			&commentData.Level,
			&commentData.CreatedAt,
			&commentData.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, &commentData)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return comments, nil
}

func (r Repo) GetCommentsByUserID(ctx context.Context, userID string) ([]*comment.Comment, error) {
	query := `
	SELECT c.id, c.content, c.post_id, c.user_id, u.username, c.parent_id, c.level, c.created_at, c.updated_at
	FROM comments c
	JOIN users u ON c.user_id = u.id
	WHERE c.user_id = ?
	ORDER BY c.created_at DESC`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments by user ID: %w", err)
	}
	defer rows.Close()

	var comments []*comment.Comment
	for rows.Next() {
		var commentData comment.Comment
		err := rows.Scan(
			&commentData.ID,
			&commentData.Content,
			&commentData.PostID,
			&commentData.UserID,
			&commentData.Username,
			&commentData.ParentID,
			&commentData.Level,
			&commentData.CreatedAt,
			&commentData.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, &commentData)
	}

	return comments, nil
}

func (r Repo) GetCommentReplies(ctx context.Context, parentID string) ([]*comment.Comment, error) {
	query := `
	SELECT c.id, c.content, c.post_id, c.user_id, u.username, c.parent_id, c.level, c.created_at, c.updated_at
	FROM comments c
	JOIN users u ON c.user_id = u.id
	WHERE c.parent_id = ?
	ORDER BY c.created_at ASC`

	rows, err := r.DB.QueryContext(ctx, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query comment replies: %w", err)
	}
	defer rows.Close()

	var comments []*comment.Comment
	for rows.Next() {
		var commentData comment.Comment
		err := rows.Scan(
			&commentData.ID,
			&commentData.Content,
			&commentData.PostID,
			&commentData.UserID,
			&commentData.Username,
			&commentData.ParentID,
			&commentData.Level,
			&commentData.CreatedAt,
			&commentData.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment reply: %w", err)
		}
		comments = append(comments, &commentData)
	}

	return comments, nil
}

func (r Repo) UpdateComment(ctx context.Context, commentData *comment.Comment) error {
	query := `
	UPDATE comments
	SET content = ?, updated_at = ?
	WHERE id = ?`

	result, err := r.DB.ExecContext(
		ctx,
		query,
		commentData.Content,
		time.Now(),
		commentData.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("comment with ID %s not found", commentData.ID)
	}

	return nil
}

func (r Repo) DeleteComment(ctx context.Context, id string) error {
	query := `DELETE FROM comments WHERE id = ?`
	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("comment with ID %s not found", id)
	}

	return nil
}

func (r Repo) GetCommentTree(ctx context.Context, postID string) ([]*comment.CommentTree, error) {
	comments, err := r.GetCommentsByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}

	return r.buildCommentTree(comments, nil), nil
}

func (r Repo) buildCommentTree(comments []*comment.Comment, parentID *string) []*comment.CommentTree {
	var tree []*comment.CommentTree

	for _, commentData := range comments {
		if (parentID == nil && commentData.ParentID == nil) ||
			(parentID != nil && commentData.ParentID != nil && *commentData.ParentID == *parentID) {

			node := &comment.CommentTree{
				Comment: commentData,
				Replies: r.buildCommentTree(comments, &commentData.ID),
			}
			tree = append(tree, node)
		}
	}

	return tree
}