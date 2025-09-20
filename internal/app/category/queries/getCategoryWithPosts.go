package queries

import (
	"context"

	"github.com/arnald/forum/internal/domain/category"
)

type GetCategoryWithPostsRequest struct {
	CategoryID string
	Limit      int
	Offset     int
}

type GetCategoryWithPostsRequestHandler interface {
	Handle(ctx context.Context, req GetCategoryWithPostsRequest) (*category.CategoryWithPosts, error)
}

type getCategoryWithPostsRequestHandler struct {
	repo category.Repository
}

func NewGetCategoryWithPostsHandler(repo category.Repository) GetCategoryWithPostsRequestHandler {
	return &getCategoryWithPostsRequestHandler{
		repo: repo,
	}
}

func (h *getCategoryWithPostsRequestHandler) Handle(ctx context.Context, req GetCategoryWithPostsRequest) (*category.CategoryWithPosts, error) {
	if req.CategoryID == "" {
		return nil, ErrEmptyID
	}

	// Set default limit if not provided
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	// Set default offset if not provided
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	return h.repo.GetCategoryWithPosts(ctx, req.CategoryID, limit, offset)
}