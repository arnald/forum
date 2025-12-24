package topics

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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

func (r Repo) GetTopicByID(ctx context.Context, topicID int, userID *string) (*topic.Topic, error) {
	query := `
	SELECT
		t.id, t.user_id, t.title, t.content, t.image_path, t.category_id, t.created_at, t.updated_at,
		u.username,
		c.name as category_name,
		c.color as category_color,
		COALESCE(vote_counts.upvotes, 0) as upvote_count,
		COALESCE(vote_counts.downvotes, 0) as downvote_count,
		COALESCE(vote_counts.score, 0) as vote_score`

	if userID != nil {
		query += `,
		user_vote.reaction_type as user_vote`
	}

	query += `
	FROM topics t
	LEFT JOIN users u ON t.user_id = u.id
	LEFT JOIN categories c ON t.category_id = c.id
	LEFT JOIN (
		SELECT
			topic_id,
			COUNT(CASE WHEN reaction_type = 1 THEN 1 END) as upvotes,
			COUNT(CASE WHEN reaction_type = -1 THEN 1 END) as downvotes,
			(COUNT(CASE WHEN reaction_type = 1 THEN 1 END) - COUNT(CASE WHEN reaction_type = -1 THEN 1 END)) as score
			FROM votes
			WHERE comment_id IS NULL
			GROUP BY topic_id
	) vote_counts ON t.id = vote_counts.topic_id`

	if userID != nil {
		query += `
		LEFT JOIN votes user_vote ON t.id = user_vote.topic_id
			AND user_vote.user_id = ?
			AND user_vote.comment_id IS NULL`
	}

	query += ` WHERE t.id = ?`

	args := make([]interface{}, 0)
	if userID != nil {
		args = append(args, *userID)
	}

	args = append(args, topicID)

	var topicResult topic.Topic
	var userVote sql.NullInt32

	scanFields := []interface{}{
		&topicResult.ID,
		&topicResult.UserID,
		&topicResult.Title,
		&topicResult.Content,
		&topicResult.ImagePath,
		&topicResult.CategoryID,
		&topicResult.CreatedAt,
		&topicResult.UpdatedAt,
		&topicResult.OwnerUsername,
		&topicResult.CategoryName,
		&topicResult.CategoryColor,
		&topicResult.UpvoteCount,
		&topicResult.DownvoteCount,
		&topicResult.VoteScore,
	}

	if userID != nil {
		scanFields = append(scanFields, &userVote)
	}

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, args...).Scan(scanFields...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("topic with ID %d not found: %w", topicID, ErrTopicNotFound)
		}
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}

	// Format Dates
	if topicResult.CreatedAt != "" {
		t, parseErr := time.Parse(time.RFC3339, topicResult.CreatedAt)
		if parseErr == nil {
			topicResult.CreatedAt = t.Format("02/01/2006")
		}
	}

	if topicResult.UpdatedAt != "" {
		t, parseErr := time.Parse(time.RFC3339, topicResult.UpdatedAt)
		if parseErr == nil {
			topicResult.UpdatedAt = t.Format("02/01/2006")
		}
	}

	if userID != nil && userVote.Valid {
		vote := int(userVote.Int32)
		topicResult.UserVote = &vote
	}

	return &topicResult, nil
}

func (r Repo) GetTotalTopicsCount(ctx context.Context, filter string, categoryID int) (int, error) {
	countQuery := `
    SELECT COUNT(*) 
    FROM topics t
    WHERE 1=1
	`

	args := make([]interface{}, 0)
	if filter != "" {
		countQuery += " AND (t.title LIKE ? OR t.content LIKE ?)"
		filterParam := "%" + filter + "%"
		args = append(args, filterParam, filterParam)
	}

	if categoryID > 0 {
		countQuery += " AND t.category_id = ?"
		args = append(args, categoryID)
	}

	var totalCount int
	err := r.DB.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return 0, fmt.Errorf("failed to get total count: %w", err)
	}

	return totalCount, nil
}

func (r Repo) GetAllTopics(ctx context.Context, page, size, categoryID int, orderBy, order, filter string, userID *string) ([]topic.Topic, error) {
	query := `
    SELECT 
        t.id, t.user_id, t.title, t.content, t.image_path, t.category_id, t.created_at, t.updated_at,
        u.username,
		c.name as category_name,
		c.color as category_color,
		COALESCE(vote_counts.upvotes, 0) as upvote_count,
		COALESCE(vote_counts.downvotes, 0) as downvote_count,
		COALESCE(vote_counts.score, 0) as vote_score
		`

	if userID != nil {
		query += `
		user_votes.reaction_type as user_vote
		`
	}

	query += `
	FROM topics t
	LEFT JOIN users u ON t.user_id = u.id
	LEFT JOIN categories c ON t.category_id = c.id
	LEFT JOIN (
		SELECT
			topic_id,
			COUNT(CASE WHEN reaction_type = 1 THEN 1 END) as upvotes,
			COUNT(CASE WHEN reaction_type = -1 THEN	1 END) as downvotes,
			(COUNT(CASE WHEN reaction_type = 1 THEN 1 END) - COUNT(CASE WHEN reaction_type = -1 THEN 1 END)) as score
			FROM votes
			WHERE comment_id IS NULL
			GROUP BY topic_id
		) vote_counts ON t.id = vote_counts.topic_id`

	if userID != nil {
		query += `
		LEFT JOIN votes user_votes ON t.id = user_votes.topic_id
			AND user_votes.user_id = ?
			AND user_votes.comment_id IS NULL`
	}

	query += ` WHERE 1=1`

	args := make([]interface{}, 0)

	if userID != nil {
		args = append(args, *userID)
	}

	if filter != "" {
		query += " AND (t.title LIKE ? OR t.content LIKE ?)"
		filterParam := "%" + filter + "%"
		args = append(args, filterParam, filterParam)
	}

	if categoryID > 0 {
		query += " AND t.category_id = ?"
		args = append(args, categoryID)
	}

	orderByClause := "t." + orderBy

	if orderBy == "vote_score" {
		orderByClause = "vote_counts.score"
	}

	query += " ORDER BY " + orderByClause + " " + order + " LIMIT ? OFFSET ?"
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
		var userVote sql.NullInt32

		scanFields := []interface{}{
			&topic.ID,
			&topic.UserID,
			&topic.Title,
			&topic.Content,
			&topic.ImagePath,
			&topic.CategoryID,
			&topic.CreatedAt,
			&topic.UpdatedAt,
			&topic.OwnerUsername,
			&topic.CategoryName,
			&topic.CategoryColor,
			&topic.UpvoteCount,
			&topic.DownvoteCount,
			&topic.VoteScore,
		}

		if userID != nil {
			scanFields = append(scanFields, &userVote)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Format dates
		if topic.CreatedAt != "" {
			t, parseErr := time.Parse(time.RFC3339, topic.CreatedAt)
			if parseErr == nil {
				topic.CreatedAt = t.Format("02/01/2006")
			}
		}

		if topic.UpdatedAt != "" {
			t, parseErr := time.Parse(time.RFC3339, topic.UpdatedAt)
			if parseErr == nil {
				topic.UpdatedAt = t.Format("02/01/2006")
			}
		}

		if userID != nil && userVote.Valid {
			vote := int(userVote.Int32)
			topic.UserVote = &vote
		}

		topic.Comments = make([]comment.Comment, 0)
		topics = append(topics, topic)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return topics, nil
}
