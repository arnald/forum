package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/arnald/forum/internal/domain/vote"
)

func (r Repo) CreateVote(ctx context.Context, voteData *vote.Vote) error {
	query := `
	INSERT INTO votes (id, user_id, target_id, target_type, vote_type, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		voteData.ID,
		voteData.UserID,
		voteData.TargetID,
		int(voteData.TargetType),
		int(voteData.VoteType),
		voteData.CreatedAt,
		voteData.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create vote: %w", err)
	}

	return nil
}

func (r Repo) UpdateVote(ctx context.Context, voteData *vote.Vote) error {
	query := `
	UPDATE votes
	SET vote_type = ?, updated_at = ?
	WHERE user_id = ? AND target_id = ? AND target_type = ?`

	result, err := r.DB.ExecContext(
		ctx,
		query,
		int(voteData.VoteType),
		voteData.UpdatedAt,
		voteData.UserID,
		voteData.TargetID,
		int(voteData.TargetType),
	)
	if err != nil {
		return fmt.Errorf("failed to update vote: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("vote not found")
	}

	return nil
}

func (r Repo) DeleteVote(ctx context.Context, userID, targetID string, targetType vote.TargetType) error {
	query := `DELETE FROM votes WHERE user_id = ? AND target_id = ? AND target_type = ?`

	result, err := r.DB.ExecContext(ctx, query, userID, targetID, int(targetType))
	if err != nil {
		return fmt.Errorf("failed to delete vote: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("vote not found")
	}

	return nil
}

func (r Repo) GetVote(ctx context.Context, userID, targetID string, targetType vote.TargetType) (*vote.Vote, error) {
	query := `
	SELECT id, user_id, target_id, target_type, vote_type, created_at, updated_at
	FROM votes
	WHERE user_id = ? AND target_id = ? AND target_type = ?`

	var voteData vote.Vote
	var targetTypeInt, voteTypeInt int

	err := r.DB.QueryRowContext(ctx, query, userID, targetID, int(targetType)).Scan(
		&voteData.ID,
		&voteData.UserID,
		&voteData.TargetID,
		&targetTypeInt,
		&voteTypeInt,
		&voteData.CreatedAt,
		&voteData.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("vote not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get vote: %w", err)
	}

	voteData.TargetType = vote.TargetType(targetTypeInt)
	voteData.VoteType = vote.VoteType(voteTypeInt)

	return &voteData, nil
}

func (r Repo) GetVoteCounts(ctx context.Context, targetID string, targetType vote.TargetType) (*vote.VoteCounts, error) {
	query := `
	SELECT
		COALESCE(SUM(CASE WHEN vote_type = 1 THEN 1 ELSE 0 END), 0) as likes,
		COALESCE(SUM(CASE WHEN vote_type = 2 THEN 1 ELSE 0 END), 0) as dislikes
	FROM votes
	WHERE target_id = ? AND target_type = ?`

	var likes, dislikes int
	err := r.DB.QueryRowContext(ctx, query, targetID, int(targetType)).Scan(&likes, &dislikes)
	if err != nil {
		return nil, fmt.Errorf("failed to get vote counts: %w", err)
	}

	return &vote.VoteCounts{
		Likes:    likes,
		Dislikes: dislikes,
		Total:    likes - dislikes,
	}, nil
}

func (r Repo) GetVoteStatus(ctx context.Context, userID, targetID string, targetType vote.TargetType) (*vote.VoteStatus, error) {
	// Get vote counts first
	counts, err := r.GetVoteCounts(ctx, targetID, targetType)
	if err != nil {
		return nil, err
	}

	// Try to get user's vote
	userVote, err := r.GetVote(ctx, userID, targetID, targetType)

	status := &vote.VoteStatus{
		VoteCounts: *counts,
		HasVoted:   false,
		VoteType:   nil,
	}

	if err == nil && userVote != nil {
		status.HasVoted = true
		status.VoteType = &userVote.VoteType
	}

	return status, nil
}

func (r Repo) GetUserVotes(ctx context.Context, userID string, targetType vote.TargetType) ([]*vote.Vote, error) {
	query := `
	SELECT id, user_id, target_id, target_type, vote_type, created_at, updated_at
	FROM votes
	WHERE user_id = ? AND target_type = ?
	ORDER BY created_at DESC`

	rows, err := r.DB.QueryContext(ctx, query, userID, int(targetType))
	if err != nil {
		return nil, fmt.Errorf("failed to query user votes: %w", err)
	}
	defer rows.Close()

	var votes []*vote.Vote
	for rows.Next() {
		var voteData vote.Vote
		var targetTypeInt, voteTypeInt int

		err := rows.Scan(
			&voteData.ID,
			&voteData.UserID,
			&voteData.TargetID,
			&targetTypeInt,
			&voteTypeInt,
			&voteData.CreatedAt,
			&voteData.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vote: %w", err)
		}

		voteData.TargetType = vote.TargetType(targetTypeInt)
		voteData.VoteType = vote.VoteType(voteTypeInt)

		votes = append(votes, &voteData)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return votes, nil
}