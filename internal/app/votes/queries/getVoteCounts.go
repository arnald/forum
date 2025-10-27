package votequeries

import (
	"context"

	"github.com/arnald/forum/internal/domain/vote"
)

type GetVoteCountsRequest struct {
	Target vote.VoteTarget
}

type getVoteCountsRequestHandler struct {
	VoteService vote.Repository
}

type GetVoteCountsRequestHandler interface {
	Handle(ctx context.Context, request GetVoteCountsRequest) (*vote.VoteCounts, error)
}

func NewGetVoteCountsRequestHandler(voteService vote.Repository) GetVoteCountsRequestHandler {
	return &getVoteCountsRequestHandler{
		VoteService: voteService,
	}
}

func (h *getVoteCountsRequestHandler) Handle(ctx context.Context, request GetVoteCountsRequest) (*vote.VoteCounts, error) {
	counts, err := h.VoteService.GetVoteCounts(ctx, request.Target)
	if err != nil {
		return nil, err
	}

	return counts, nil
}
