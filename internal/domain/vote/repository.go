package vote

import "context"

type Repository interface {
	CastVote(ctx context.Context, userID string, target VoteTarget, reactionType int) error
	DeleteVote(ctx context.Context, voteID int, userID string) error
}
