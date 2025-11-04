package vote

import "context"

type Repository interface {
	CastVote(ctx context.Context, userID string, target Target, reactionType int) error
	DeleteVote(ctx context.Context, voteID int, userID string) error
	GetCounts(ctx context.Context, target Target) (*Counts, error)
}
