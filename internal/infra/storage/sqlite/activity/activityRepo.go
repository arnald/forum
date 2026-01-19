package activities

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/arnald/forum/internal/domain/activity"
)

type Repo struct {
	DB *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{DB: db}
}

func (r *Repo) GetUserActivity(ctx context.Context, userID string) (*activity.Activity, error) {
	act := &activity.Activity{}

	createdTopics, err := r.getCreatedTopics(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created topics: %w", err)
	}
	act.CreatedTopics = createdTopics

	likedTopics, err := r.getVotedTopics(ctx, userID, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get liked topics: %w", err)
	}
	act.LikedTopics = likedTopics

	dislikedTopics, err := r.getVotedTopics(ctx, userID, -1)
	if err != nil {
		return nil, fmt.Errorf("failed to get disliked topics: %w", err)
	}
	act.DislikedTopics = dislikedTopics

	likedComments, err := r.getVotedComments(ctx, userID, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get liked comments: %w", err)
	}
	act.LikedComments = likedComments

	dislikedComments, err := r.getVotedComments(ctx, userID, -1)
	if err != nil {
		return nil, fmt.Errorf("failed to get disliked comments: %w", err)
	}
	act.DislikedComments = dislikedComments

	userComments, err := r.getUserComments(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user comments: %w", err)
	}
	act.UserComments = userComments

	return act, nil
}

func (r *Repo) getCreatedTopics(ctx context.Context, userID string) ([]activity.TopicActivity, error) {
	query := `
        SELECT id, title, created_at
        FROM topics
        WHERE user_id = ?
        ORDER BY created_at DESC
        LIMIT 50`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	topics := make([]activity.TopicActivity, 0)
	for rows.Next() {
		var topic activity.TopicActivity
		var createdAt string
		rowsErr := rows.Scan(&topic.ID, &topic.Title, &createdAt)
		if rowsErr != nil {
			return nil, rowsErr
		}

		t, parseErr := time.Parse(time.RFC3339, createdAt)
		if parseErr == nil {
			topic.CreatedAt = t.Format("Jan 2, 2006 3:04 PM")
		} else {
			topic.CreatedAt = createdAt
		}

		topics = append(topics, topic)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return topics, nil
}

func (r *Repo) getVotedTopics(ctx context.Context, userID string, reactionType int) ([]activity.TopicActivity, error) {
	query := `
        SELECT t.id, t.title, v.created_at
        FROM topics t
        INNER JOIN votes v ON t.id = v.topic_id
        WHERE v.user_id = ? 
        AND v.reaction_type = ? 
        AND v.comment_id IS NULL
        ORDER BY v.created_at DESC
        LIMIT 50`

	rows, err := r.DB.QueryContext(ctx, query, userID, reactionType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	topics := make([]activity.TopicActivity, 0)
	for rows.Next() {
		var topic activity.TopicActivity
		var createdAt string
		rowsErr := rows.Scan(&topic.ID, &topic.Title, &createdAt)
		if rowsErr != nil {
			return nil, rowsErr
		}

		t, parseErr := time.Parse(time.RFC3339, createdAt)
		if parseErr == nil {
			topic.CreatedAt = t.Format("Jan 2, 2006 3:04 PM")
		} else {
			topic.CreatedAt = createdAt
		}

		topics = append(topics, topic)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return topics, nil
}

func (r *Repo) getVotedComments(ctx context.Context, userID string, reactionType int) ([]activity.CommentVoteActivity, error) {
	query := `
        SELECT c.id, c.topic_id, t.title, v.created_at
        FROM comments c
        INNER JOIN votes v ON c.id = v.comment_id
        INNER JOIN topics t ON c.topic_id = t.id
        WHERE v.user_id = ? 
        AND v.reaction_type = ?
        ORDER BY v.created_at DESC
        LIMIT 50`

	rows, err := r.DB.QueryContext(ctx, query, userID, reactionType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]activity.CommentVoteActivity, 0)
	for rows.Next() {
		var comment activity.CommentVoteActivity
		var createdAt string
		rowsErr := rows.Scan(&comment.CommentID, &comment.TopicID, &comment.TopicTitle, &createdAt)
		if rowsErr != nil {
			return nil, rowsErr
		}

		t, parseErr := time.Parse(time.RFC3339, createdAt)
		if parseErr == nil {
			comment.CreatedAt = t.Format("Jan 2, 2006 3:04 PM")
		} else {
			comment.CreatedAt = createdAt
		}

		comments = append(comments, comment)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *Repo) getUserComments(ctx context.Context, userID string) ([]activity.CommentActivity, error) {
	query := `
        SELECT c.id, c.content, c.topic_id, t.title, c.created_at
        FROM comments c
        INNER JOIN topics t ON c.topic_id = t.id
        WHERE c.user_id = ?
        ORDER BY c.created_at DESC
        LIMIT 50`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]activity.CommentActivity, 0)
	for rows.Next() {
		var comment activity.CommentActivity
		var createdAt string
		rowsErr := rows.Scan(&comment.ID, &comment.Content, &comment.TopicID, &comment.TopicTitle, &createdAt)
		if rowsErr != nil {
			return nil, rowsErr
		}

		// Parse and format date
		t, parseErr := time.Parse(time.RFC3339, createdAt)
		if parseErr == nil {
			comment.CreatedAt = t.Format("Jan 2, 2006 3:04 PM")
		} else {
			comment.CreatedAt = createdAt
		}

		comments = append(comments, comment)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return comments, nil
}
