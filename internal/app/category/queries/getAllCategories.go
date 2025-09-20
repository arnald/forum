package queries

import (
	"context"

	"github.com/arnald/forum/internal/domain/category"
)

type GetAllCategoriesRequest struct {
	// No parameters needed for getting all categories
}

type GetAllCategoriesRequestHandler interface {
	Handle(ctx context.Context, req GetAllCategoriesRequest) ([]*category.Category, error)
}

type getAllCategoriesRequestHandler struct {
	repo category.Repository
}

func NewGetAllCategoriesHandler(repo category.Repository) GetAllCategoriesRequestHandler {
	return &getAllCategoriesRequestHandler{
		repo: repo,
	}
}

func (h *getAllCategoriesRequestHandler) Handle(ctx context.Context, req GetAllCategoriesRequest) ([]*category.Category, error) {
	return h.repo.GetAllCategories(ctx)
}