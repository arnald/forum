package categorycommands

import (
	"context"

	"github.com/arnald/forum/internal/domain/category"
)

type CreateCategoryRequest struct {
	Name        string
	Description string
	CreatedBy   string
}

type CreateCategoryRequestHandler interface {
	Handle(ctx context.Context, req CreateCategoryRequest) error
}

type createCategoryRequestHandler struct {
	repo category.Repository
}

func NewCreateCategoryHandler(repo category.Repository) CreateCategoryRequestHandler {
	return &createCategoryRequestHandler{
		repo: repo,
	}
}

func (h *createCategoryRequestHandler) Handle(ctx context.Context, req CreateCategoryRequest) error {
	category := &category.Category{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   req.CreatedBy,
	}

	err := h.repo.CreateCategory(ctx, category)
	if err != nil {
		return err
	}
	return nil
}
