package categorycommands

import (
	"context"

	"github.com/arnald/forum/internal/domain/category"
)

type DeleteCategoryRequest struct {
	UserID     string
	CategoryID int
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
	err := h.repo.DeleteCategory(ctx, req.CategoryID, req.UserID)
	if err != nil {
		return err
	}
	return nil
}
