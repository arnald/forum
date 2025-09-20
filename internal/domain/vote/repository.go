package vote

import (
	"context"
)

type Repository interface {
	CreateVote(ctx context.Context, vote *Vote) error
	UpdateVote(ctx context.Context, vote *Vote) error
	DeleteVote(ctx context.Context, userID, targetID string, targetType TargetType) error
	GetVote(ctx context.Context, userID, targetID string, targetType TargetType) (*Vote, error)
	GetVoteCounts(ctx context.Context, targetID string, targetType TargetType) (*VoteCounts, error)
	GetVoteStatus(ctx context.Context, userID, targetID string, targetType TargetType) (*VoteStatus, error)
	GetUserVotes(ctx context.Context, userID string, targetType TargetType) ([]*Vote, error)
}