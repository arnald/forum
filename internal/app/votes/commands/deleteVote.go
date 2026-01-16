package votecommands

import (
	"context"

	"github.com/arnald/forum/internal/domain/vote"
)

type DeleteVoteRequest struct {
	UserID    string
	TopicID   *int
	CommentID *int
}

type deleteVoteRequestHandler struct {
	VoteRepo vote.Repository
}

type DeleteVoteRequestHandler interface {
	Handle(ctx context.Context, req DeleteVoteRequest) error
}

func NewDeleteVoteHandler(repo vote.Repository) DeleteVoteRequestHandler {
	return &deleteVoteRequestHandler{
		VoteRepo: repo,
	}
}

func (h *deleteVoteRequestHandler) Handle(ctx context.Context, req DeleteVoteRequest) error {
	err := h.VoteRepo.DeleteVote(ctx, req.UserID, req.TopicID, req.CommentID)
	if err != nil {
		return err
	}

	return nil
}
