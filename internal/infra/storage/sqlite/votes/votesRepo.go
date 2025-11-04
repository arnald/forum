package votes

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/arnald/forum/internal/domain/vote"
)

type Repo struct {
	DB *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{DB: db}
}

func (r *Repo) CastVote(ctx context.Context, userID string, target vote.Target, reactionType int) error {
	var query string
	var args []interface{}

	if target.CommentID != nil {
		query = `
		INSERT INTO votes (user_id, topic_id, comment_id, reaction_type) 
		VALUES (?, NULL, ?, ?) 
		ON CONFLICT (user_id, comment_id) DO UPDATE SET 
			reaction_type = EXCLUDED.reaction_type,
			created_at = CURRENT_TIMESTAMP`
		args = []interface{}{userID, *target.CommentID, reactionType}
	} else {
		query = `
		INSERT INTO votes (user_id, topic_id, comment_id, reaction_type) 
		VALUES (?, ?, NULL, ?) 
		ON CONFLICT (user_id, topic_id) DO UPDATE SET 
			reaction_type = EXCLUDED.reaction_type,
			created_at = CURRENT_TIMESTAMP`
		args = []interface{}{userID, target.TopicID, reactionType}
	}

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query for casting vote: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return fmt.Errorf("failed to cast vote: %w", err)
	}

	return nil
}

func (r *Repo) DeleteVote(ctx context.Context, voteID int, userID string) error {
	query := `
	DELETE FROM votes
	WHERE id = ? AND user_id = ?
	`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, voteID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete vote: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return ErrVoteNotFound
	}

	return nil
}

func (r *Repo) GetCounts(ctx context.Context, target vote.Target) (*vote.Counts, error) {
	var query string
	var args []interface{}

	if target.CommentID == nil {
		query = `
		SELECT
			COUNT(CASE WHEN reaction_type = 1 THEN 1 END) as upvotes,
			COUNT(CASE WHEN reaction_type = -1 THEN 1 END) as downvotes
		FROM votes
		WHERE topic_id = ? AND comment_id IS NULL`
		args = []interface{}{target.TopicID}
	} else {
		query = `
		SELECT
			COUNT(CASE WHEN reaction_type = 1 THEN 1 END) as upvotes,
			COUNT(CASE WHEN reaction_type = -1 THEN 1 END) as downvotes
		FROM votes
		WHERE comment_id = ? AND topic_id IS NULL`
		args = []interface{}{target.CommentID}
	}

	var counts vote.Counts
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, args...).Scan(
		&counts.Upvotes,
		&counts.DownVotes,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get vote counts: %w", err)
	}

	counts.Score = counts.Upvotes - counts.DownVotes

	return &counts, nil
}
