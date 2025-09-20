package queries

import (
	"context"

	"github.com/arnald/forum/internal/domain/vote"
)

type GetUserVotesRequest struct {
	UserID     string
	TargetType vote.TargetType
}

type GetUserVotesRequestHandler interface {
	Handle(ctx context.Context, req GetUserVotesRequest) ([]*vote.Vote, error)
}

type getUserVotesRequestHandler struct {
	repo vote.Repository
}

func NewGetUserVotesHandler(repo vote.Repository) GetUserVotesRequestHandler {
	return &getUserVotesRequestHandler{
		repo: repo,
	}
}

func (h *getUserVotesRequestHandler) Handle(ctx context.Context, req GetUserVotesRequest) ([]*vote.Vote, error) {
	if req.UserID == "" {
		return nil, ErrEmptyUserID
	}

	return h.repo.GetUserVotes(ctx, req.UserID, req.TargetType)
}