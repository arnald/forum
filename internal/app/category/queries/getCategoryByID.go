package queries

import (
	"context"

	"github.com/arnald/forum/internal/domain/category"
)

type GetCategoryByIDRequest struct {
	ID string
}

type GetCategoryByIDRequestHandler interface {
	Handle(ctx context.Context, req GetCategoryByIDRequest) (*category.Category, error)
}

type getCategoryByIDRequestHandler struct {
	repo category.Repository
}

func NewGetCategoryByIDHandler(repo category.Repository) GetCategoryByIDRequestHandler {
	return &getCategoryByIDRequestHandler{
		repo: repo,
	}
}

func (h *getCategoryByIDRequestHandler) Handle(ctx context.Context, req GetCategoryByIDRequest) (*category.Category, error) {
	if req.ID == "" {
		return nil, ErrEmptyID
	}

	return h.repo.GetCategoryByID(ctx, req.ID)
}