package category

import "context"

type Repository interface {
	CreateCategory(ctx context.Context, category *Category) error
	DeleteCategory(ctx context.Context, id int, userID string) error
	UpdateCategory(ctx context.Context, category *Category) error
	GetCategoryByID(ctx context.Context, id int) (*Category, error)
	GetAllCategories(ctx context.Context) ([]*Category, error)
}
