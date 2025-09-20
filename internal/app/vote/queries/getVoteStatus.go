package queries

import (
	"context"

	"github.com/arnald/forum/internal/domain/vote"
)

type GetVoteStatusRequest struct {
	UserID     string
	TargetID   string
	TargetType vote.TargetType
}

type GetVoteStatusRequestHandler interface {
	Handle(ctx context.Context, req GetVoteStatusRequest) (*vote.VoteStatus, error)
}

type getVoteStatusRequestHandler struct {
	repo vote.Repository
}

func NewGetVoteStatusHandler(repo vote.Repository) GetVoteStatusRequestHandler {
	return &getVoteStatusRequestHandler{
		repo: repo,
	}
}

func (h *getVoteStatusRequestHandler) Handle(ctx context.Context, req GetVoteStatusRequest) (*vote.VoteStatus, error) {
	if req.TargetID == "" {
		return nil, ErrEmptyTargetID
	}

	return h.repo.GetVoteStatus(ctx, req.UserID, req.TargetID, req.TargetType)
}