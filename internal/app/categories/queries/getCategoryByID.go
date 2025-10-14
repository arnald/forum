package categoryqueries

import (
	"context"

	"github.com/arnald/forum/internal/domain/category"
)

type GetCategoryByIDRequest struct {
	ID int
}

type GetCategoryByIDHandler interface {
	Handle(ctx context.Context, req GetCategoryByIDRequest) (*category.Category, error)
}

type getCategoryByIDHandler struct {
	repo category.Repository
}

func NewGetCategoryByIDHandler(repo category.Repository) GetCategoryByIDHandler {
	return &getCategoryByIDHandler{repo: repo}
}

func (h *getCategoryByIDHandler) Handle(ctx context.Context, req GetCategoryByIDRequest) (*category.Category, error) {
	category, err := h.repo.GetCategoryByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return category, nil
}
