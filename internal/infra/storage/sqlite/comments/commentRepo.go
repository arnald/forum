package comments

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/arnald/forum/internal/domain/comment"
)

type Repo struct {
	DB *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		DB: db,
	}
}

func (r *Repo) CreateComment(ctx context.Context, comment *comment.Comment) error {
	query := `
	INSERT INTO comments (user_id, topic_id, content)
	VALUES (?, ?, ?)`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		comment.UserID,
		comment.TopicID,
		comment.Content,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return fmt.Errorf("user or topic not found: %w", err)
		default:
			return fmt.Errorf("failed to create comment: %w", err)
		}
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	comment.ID = int(id)
	return nil
}

func (r *Repo) UpdateComment(ctx context.Context, comment *comment.Comment) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				err = fmt.Errorf("transaction rollback failed: %w (original error: %w)", rollbackErr, err)
			}
			return
		}
		commitErr := tx.Commit()
		if commitErr != nil {
			err = fmt.Errorf("transaction commit failed: %w", commitErr)
		}
	}()

	query := `
	UPDATE comments 
	SET content = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ? AND user_id = ?`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		comment.Content,
		comment.ID,
		comment.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to execute update: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comment with ID %d %w", comment.ID, ErrFailedToUpdate)
	}

	return nil
}

func (r *Repo) DeleteComment(ctx context.Context, userID string, commentID int) error {
	query := `
	DELETE FROM comments
	WHERE id = ? AND user_id = ?`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, commentID, userID)
	if err != nil {
		return fmt.Errorf("failed to execute delete statement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comment with ID %d not found or user not authorized: %w", commentID, ErrCommentNotFound)
	}

	return nil
}

func (r *Repo) GetCommentByID(ctx context.Context, commentID int) (*comment.Comment, error) {
	query := `
	SELECT 
		c.id, c.user_id, c.topic_id, c.content, c.created_at, c.updated_at, u.username
	FROM comments c
	LEFT JOIN users u ON c.user_id = u.id
	WHERE c.id = ?`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	comment := &comment.Comment{}
	err = stmt.QueryRowContext(ctx, commentID).Scan(
		&comment.ID,
		&comment.UserID,
		&comment.TopicID,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.OwnerUsername,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("comment with ID %d not found: %w", commentID, ErrCommentNotFound)
		}
		return nil, fmt.Errorf("failed to query comment: %w", err)
	}

	return comment, nil
}

func (r *Repo) GetCommentsByTopicID(ctx context.Context, topicID int) ([]comment.Comment, error) {
	query := `
	SELECT 
		c.id, c.user_id, c.topic_id, c.content, c.created_at, c.updated_at, u.username
	FROM comments c
	LEFT JOIN users u ON c.user_id = u.id
	WHERE c.topic_id = ?
	ORDER BY c.created_at ASC`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, topicID)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	comments := make([]comment.Comment, 0)
	for rows.Next() {
		var c comment.Comment
		err = rows.Scan(
			&c.ID,
			&c.UserID,
			&c.TopicID,
			&c.Content,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.OwnerUsername,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		comments = append(comments, c)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return comments, nil
}
