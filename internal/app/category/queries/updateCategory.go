package queries

import (
	"context"
	"strings"

	"github.com/arnald/forum/internal/domain/category"
)

type UpdateCategoryRequest struct {
	ID          string
	Name        string
	Description string
}

type UpdateCategoryRequestHandler interface {
	Handle(ctx context.Context, req UpdateCategoryRequest) (*category.Category, error)
}

type updateCategoryRequestHandler struct {
	repo category.Repository
}

func NewUpdateCategoryHandler(repo category.Repository) UpdateCategoryRequestHandler {
	return &updateCategoryRequestHandler{
		repo: repo,
	}
}

func (h *updateCategoryRequestHandler) Handle(ctx context.Context, req UpdateCategoryRequest) (*category.Category, error) {
	if req.ID == "" {
		return nil, ErrEmptyID
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrEmptyName
	}

	// Get existing category
	existingCategory, err := h.repo.GetCategoryByID(ctx, req.ID)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	// Check if name is being changed and if new name already exists
	if existingCategory.Name != name {
		existingByName, err := h.repo.GetCategoryByName(ctx, name)
		if err == nil && existingByName != nil && existingByName.ID != req.ID {
			return nil, ErrCategoryAlreadyExists
		}
	}

	updatedCategory := &category.Category{
		ID:          req.ID,
		Name:        name,
		Description: strings.TrimSpace(req.Description),
		CreatedAt:   existingCategory.CreatedAt,
		PostCount:   existingCategory.PostCount,
	}

	err = h.repo.UpdateCategory(ctx, updatedCategory)
	if err != nil {
		return nil, err
	}

	return h.repo.GetCategoryByID(ctx, req.ID)
}