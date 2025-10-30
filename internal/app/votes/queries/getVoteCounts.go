package votequeries

import (
	"context"

	"github.com/arnald/forum/internal/domain/vote"
)

type GetCountsRequest struct {
	Target vote.Target
}

type getCountsRequestHandler struct {
	VoteService vote.Repository
}

type GetCountsRequestHandler interface {
	Handle(ctx context.Context, request GetCountsRequest) (*vote.Counts, error)
}

func NewGetCountsRequestHandler(voteService vote.Repository) GetCountsRequestHandler {
	return &getCountsRequestHandler{
		VoteService: voteService,
	}
}

func (h *getCountsRequestHandler) Handle(ctx context.Context, request GetCountsRequest) (*vote.Counts, error) {
	counts, err := h.VoteService.GetCounts(ctx, request.Target)
	if err != nil {
		return nil, err
	}

	return counts, nil
}
