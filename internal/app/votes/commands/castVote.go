package votecommands

import (
	"context"

	"github.com/arnald/forum/internal/domain/vote"
)

type CastVoteRequest struct {
	Target       vote.Target `json:"target"`
	UserID       string      `json:"userId"`
	ReactionType int         `json:"reactionType"`
}

type castVoteRequestHandler struct {
	VoteRepo vote.Repository
}

type CastVoteRequestHandler interface {
	Handle(ctx context.Context, req CastVoteRequest) error
}

func NewCastVoteHandler(voteRepo vote.Repository) CastVoteRequestHandler {
	return &castVoteRequestHandler{
		VoteRepo: voteRepo,
	}
}

func (h *castVoteRequestHandler) Handle(ctx context.Context, req CastVoteRequest) error {
	err := h.VoteRepo.CastVote(ctx, req.UserID, req.Target, req.ReactionType)
	if err != nil {
		return err
	}

	return nil
}
