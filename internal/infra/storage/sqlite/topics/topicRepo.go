package topics

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/arnald/forum/internal/domain/comment"
	"github.com/arnald/forum/internal/domain/topic"
)

type Repo struct {
	DB *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		DB: db,
	}
}

func (r Repo) CreateTopic(ctx context.Context, topic *topic.Topic) error {
	query := `
	INSERT INTO topics (user_id, title, content, image_path, category_id)
	VALUES (?, ?, ?, ?, ?)`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		topic.UserID,
		topic.Title,
		topic.Content,
		topic.ImagePath,
		topic.CategoryID,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return fmt.Errorf("user with ID %s not found: %w", topic.UserID, ErrUserNotFound)
		default:
			return fmt.Errorf("failed to create topic: %w", err)
		}
	}

	return nil
}

func (r Repo) UpdateTopic(ctx context.Context, topic *topic.Topic) error {
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
		err = tx.Commit()
		if err != nil {
			err = fmt.Errorf("transaction commit failed: %w", err)
		}
	}()

	query := `
	UPDATE topics 
	SET title = ?, content = ?, image_path = ?, category_id = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ? AND user_id = ?`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		topic.Title,
		topic.Content,
		topic.ImagePath,
		topic.CategoryID,
		topic.ID,
		topic.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to execute update: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("topic with ID %d not found or user not authorized: %w", topic.ID, ErrTopicNotFound)
	}

	return nil
}

func (r Repo) DeleteTopic(ctx context.Context, userID string, topicID int) error {
	query := `
	DELETE FROM topics
	WHERE id = ? AND user_id = ?`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, topicID, userID)
	if err != nil {
		return fmt.Errorf("failed to execute delete statement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("topic with ID %d not found or user not authorized: %w", topicID, ErrTopicNotFound)
	}

	return nil
}

func (r Repo) GetTopicByID(ctx context.Context, topicID int) (*topic.Topic, error) {
	query := `
	SELECT 
		t.id, t.user_id, u.username, t.title, t.content, t.image_path, t.category_id, t.created_at, t.updated_at,
		c.id, c.user_id, c.content, c.created_at, c.updated_at, u.username
	FROM topics t
	LEFT JOIN comments c ON t.id = c.topic_id
	LEFT JOIN users u ON c.user_id = u.id
	WHERE t.id = ?
	`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, topicID)
	if err != nil {
		return nil, fmt.Errorf("failed to query topic: %w", err)
	}
	defer rows.Close()

	topic := &topic.Topic{}
	commentsList := make([]comment.Comment, 0)
	found := false

	for rows.Next() {
		found = true

		var commentID sql.NullInt64
		var commentUserID, commentContent, commentUsername sql.NullString
		var commentCreatedAt, commentUpdatedAt sql.NullTime

		err = rows.Scan(
			&topic.ID,
			&topic.UserID,
			&topic.OwnerUsername,
			&topic.Title,
			&topic.Content,
			&topic.ImagePath,
			&topic.CategoryID,
			&topic.CreatedAt,
			&topic.UpdatedAt,
			&commentID,
			&commentUserID,
			&commentContent,
			&commentCreatedAt,
			&commentUpdatedAt,
			&commentUsername,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if commentID.Valid {
			comment := comment.Comment{
				ID:        int(commentID.Int64),
				UserID:    commentUserID.String,
				TopicID:   topicID,
				Content:   commentContent.String,
				CreatedAt: commentCreatedAt.Time,
				UpdatedAt: commentUpdatedAt.Time,
				Username:  commentUsername.String,
			}
			commentsList = append(commentsList, comment)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	if !found {
		return nil, fmt.Errorf("topic with ID %d not found: %w", topicID, ErrTopicNotFound)
	}

	topic.Comments = commentsList
	return topic, nil
}

func (r Repo) GetTotalTopicsCount(ctx context.Context, filter string) (int, error) {
	countQuery := `
    SELECT COUNT(*) 
    FROM topics t
    WHERE 1=1`

	args := make([]interface{}, 0)
	if filter != "" {
		countQuery += " AND (t.title LIKE ? OR t.content LIKE ?)"
		filterParam := "%" + filter + "%"
		args = append(args, filterParam, filterParam)
	}

	var totalCount int
	err := r.DB.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return 0, fmt.Errorf("failed to get total count: %w", err)
	}

	return totalCount, nil
}

func (r Repo) GetAllTopics(ctx context.Context, page, size int, orderBy, filter string) ([]topic.Topic, error) {
	query := `
    SELECT 
        t.id, t.user_id, t.title, t.content, t.image_path, t.category_id, t.created_at, t.updated_at,
        u.username
    FROM topics t
    LEFT JOIN users u ON t.user_id = u.id
    WHERE 1=1`

	args := make([]interface{}, 0)
	if filter != "" {
		query += " AND (t.title LIKE ? OR t.content LIKE ?)"
		filterParam := "%" + filter + "%"
		args = append(args, filterParam, filterParam)
	}

	query += " ORDER BY t." + orderBy + " LIMIT ? OFFSET ?"
	offset := (page - 1) * size
	args = append(args, size, offset)

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query topics: %w", err)
	}
	defer rows.Close()

	topics := make([]topic.Topic, 0)
	for rows.Next() {
		var topic topic.Topic
		err = rows.Scan(
			&topic.ID,
			&topic.UserID,
			&topic.Title,
			&topic.Content,
			&topic.ImagePath,
			&topic.CategoryID,
			&topic.CreatedAt,
			&topic.UpdatedAt,
			&topic.OwnerUsername,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		topics = append(topics, topic)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return topics, nil
}
