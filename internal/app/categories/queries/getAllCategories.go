package categoryqueries

import (
	"context"

	"github.com/arnald/forum/internal/domain/category"
)

type GetAllCategoriesRequest struct {
	OrderBy string `json:"orderBy"`
	Order   string `json:"order"`
	Filter  string `json:"filter"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
	Offset  int    `json:"offset"`
}

type GetAllCategoriesRequestHandler interface {
	Handle(ctx context.Context, req GetAllCategoriesRequest) ([]category.Category, int, error)
}

type getAllCategoriesRequestHandler struct {
	repo category.Repository
}

func NewGetAllCategoriesHandler(repo category.Repository) GetAllCategoriesRequestHandler {
	return getAllCategoriesRequestHandler{
		repo: repo,
	}
}

func (h getAllCategoriesRequestHandler) Handle(ctx context.Context, req GetAllCategoriesRequest) ([]category.Category, int, error) {
	categories, err := h.repo.GetAllCategories(
		ctx,
		req.Page,
		req.Size,
		req.OrderBy,
		req.Order,
		req.Filter,
	)
	if err != nil {
		return nil, 0, err
	}

	count, err := h.repo.GetTotalCategoriesCount(ctx, req.Filter)
	if err != nil {
		return nil, 0, err
	}
	return categories, count, nil
}
