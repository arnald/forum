package vote

import "context"

type Repository interface {
	CastVote(ctx context.Context, userID string, target Target, reactionType int) error
	DeleteVote(ctx context.Context, userID string, topicID *int, coommentID *int) error
	GetCounts(ctx context.Context, target Target) (*Counts, error)
}
