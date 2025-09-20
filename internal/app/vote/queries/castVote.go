package queries

import (
	"context"
	"time"

	"github.com/arnald/forum/internal/domain/vote"
	"github.com/arnald/forum/internal/pkg/uuid"
)

type CastVoteRequest struct {
	UserID     string
	TargetID   string
	TargetType vote.TargetType
	VoteType   vote.VoteType
}

type CastVoteRequestHandler interface {
	Handle(ctx context.Context, req CastVoteRequest) (*vote.VoteStatus, error)
}

type castVoteRequestHandler struct {
	repo         vote.Repository
	uuidProvider uuid.Provider
}

func NewCastVoteHandler(repo vote.Repository, uuidProvider uuid.Provider) CastVoteRequestHandler {
	return &castVoteRequestHandler{
		repo:         repo,
		uuidProvider: uuidProvider,
	}
}

func (h *castVoteRequestHandler) Handle(ctx context.Context, req CastVoteRequest) (*vote.VoteStatus, error) {
	if req.UserID == "" {
		return nil, ErrEmptyUserID
	}
	if req.TargetID == "" {
		return nil, ErrEmptyTargetID
	}

	// Check if user already voted
	existingVote, err := h.repo.GetVote(ctx, req.UserID, req.TargetID, req.TargetType)

	if err != nil && existingVote == nil {
		// No existing vote, create new one
		newVote := &vote.Vote{
			ID:         h.uuidProvider.NewUUID(),
			UserID:     req.UserID,
			TargetID:   req.TargetID,
			TargetType: req.TargetType,
			VoteType:   req.VoteType,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		err = h.repo.CreateVote(ctx, newVote)
		if err != nil {
			return nil, err
		}
	} else if existingVote != nil {
		// User already voted
		if existingVote.VoteType == req.VoteType {
			// Same vote type, remove the vote (toggle off)
			err = h.repo.DeleteVote(ctx, req.UserID, req.TargetID, req.TargetType)
			if err != nil {
				return nil, err
			}
		} else {
			// Different vote type, update the vote
			existingVote.VoteType = req.VoteType
			existingVote.UpdatedAt = time.Now()
			err = h.repo.UpdateVote(ctx, existingVote)
			if err != nil {
				return nil, err
			}
		}
	}

	// Return updated vote status
	return h.repo.GetVoteStatus(ctx, req.UserID, req.TargetID, req.TargetType)
}