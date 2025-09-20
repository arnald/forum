package queries

import (
	"context"

	"github.com/arnald/forum/internal/domain/category"
)

type DeleteCategoryRequest struct {
	ID string
}

type DeleteCategoryRequestHandler interface {
	Handle(ctx context.Context, req DeleteCategoryRequest) error
}

type deleteCategoryRequestHandler struct {
	repo category.Repository
}

func NewDeleteCategoryHandler(repo category.Repository) DeleteCategoryRequestHandler {
	return &deleteCategoryRequestHandler{
		repo: repo,
	}
}

func (h *deleteCategoryRequestHandler) Handle(ctx context.Context, req DeleteCategoryRequest) error {
	if req.ID == "" {
		return ErrEmptyID
	}

	// Check if category exists
	_, err := h.repo.GetCategoryByID(ctx, req.ID)
	if err != nil {
		return ErrCategoryNotFound
	}

	return h.repo.DeleteCategory(ctx, req.ID)
}