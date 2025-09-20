package queries

import (
	"context"
	"strings"
	"time"

	"github.com/arnald/forum/internal/domain/category"
	"github.com/arnald/forum/internal/pkg/uuid"
)

type CreateCategoryRequest struct {
	Name        string
	Description string
}

type CreateCategoryRequestHandler interface {
	Handle(ctx context.Context, req CreateCategoryRequest) (*category.Category, error)
}

type createCategoryRequestHandler struct {
	repo         category.Repository
	uuidProvider uuid.Provider
}

func NewCreateCategoryHandler(repo category.Repository, uuidProvider uuid.Provider) CreateCategoryRequestHandler {
	return &createCategoryRequestHandler{
		repo:         repo,
		uuidProvider: uuidProvider,
	}
}

func (h *createCategoryRequestHandler) Handle(ctx context.Context, req CreateCategoryRequest) (*category.Category, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrEmptyName
	}

	// Check if category with this name already exists
	existingCategory, err := h.repo.GetCategoryByName(ctx, name)
	if err == nil && existingCategory != nil {
		return nil, ErrCategoryAlreadyExists
	}

	newCategory := &category.Category{
		ID:          h.uuidProvider.NewUUID(),
		Name:        name,
		Description: strings.TrimSpace(req.Description),
		CreatedAt:   time.Now(),
		PostCount:   0,
	}

	err = h.repo.CreateCategory(ctx, newCategory)
	if err != nil {
		return nil, err
	}

	return newCategory, nil
}