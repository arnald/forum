package queries

import (
	"context"

	"github.com/arnald/forum/internal/domain/post"
)

type GetPostsByCategoryRequest struct {
	CategoryID   string
	CategoryName string
}

type GetPostsByCategoryRequestHandler interface {
	Handle(ctx context.Context, req GetPostsByCategoryRequest) ([]*post.Post, error)
}

type getPostsByCategoryRequestHandler struct {
	repo post.Repository
}

func NewGetPostsByCategoryHandler(repo post.Repository) GetPostsByCategoryRequestHandler {
	return &getPostsByCategoryRequestHandler{
		repo: repo,
	}
}

func (h *getPostsByCategoryRequestHandler) Handle(ctx context.Context, req GetPostsByCategoryRequest) ([]*post.Post, error) {
	if req.CategoryName == "" && req.CategoryID == "" {
		return nil, ErrEmptyCategory
	}

	// If we have category ID, use it (more reliable than name)
	if req.CategoryID != "" {
		return h.repo.GetPostsByCategoryID(ctx, req.CategoryID)
	}

	// If we only have category name, use name-based filtering
	if req.CategoryName != "" {
		return h.repo.GetPostsByCategory(ctx, req.CategoryName)
	}

	return nil, ErrEmptyCategory
}