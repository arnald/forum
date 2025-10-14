package categorycommands

import (
	"context"

	"github.com/arnald/forum/internal/domain/category"
)

type UpdateCategoryRequest struct {
	Name        string
	Description string
	ID          int
}

type UpdateCategoryRequestHandler interface {
	Handle(ctx context.Context, req UpdateCategoryRequest) error
}

type updateCategoryRequestHandler struct {
	repo category.Repository
}

func NewUpdateCategoryHandler(repo category.Repository) UpdateCategoryRequestHandler {
	return &updateCategoryRequestHandler{
		repo: repo,
	}
}

func (h *updateCategoryRequestHandler) Handle(ctx context.Context, req UpdateCategoryRequest) error {
	category := &category.Category{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
	}
	err := h.repo.UpdateCategory(ctx, category)
	if err != nil {
		return err
	}
	return nil
}
