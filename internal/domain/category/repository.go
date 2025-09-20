package category

import (
	"context"
)

type Repository interface {
	CreateCategory(ctx context.Context, category *Category) error
	GetCategoryByID(ctx context.Context, id string) (*Category, error)
	GetCategoryByName(ctx context.Context, name string) (*Category, error)
	GetAllCategories(ctx context.Context) ([]*Category, error)
	UpdateCategory(ctx context.Context, category *Category) error
	DeleteCategory(ctx context.Context, id string) error
	GetCategoryWithPosts(ctx context.Context, categoryID string, limit int, offset int) (*CategoryWithPosts, error)
	GetPostsCount(ctx context.Context, categoryID string) (int, error)
}