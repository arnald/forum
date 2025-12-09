package category

import (
	"context"
)

type Repository interface {
	CreateCategory(ctx context.Context, category *Category) error
	DeleteCategory(ctx context.Context, id int, userID string) error
	UpdateCategory(ctx context.Context, category *Category) error
	GetCategoryByID(ctx context.Context, id int) (*Category, error)
	GetAllCategories(ctx context.Context, page, size int, orderBy, order, filter string) ([]Category, error)
	PopulateCategoriesWithTopics(ctx context.Context, categories []Category) ([]Category, error)
	GetTotalCategoriesCount(ctx context.Context, filter string) (int, error)
	GetAllCategorieNamesAndIDs(ctx context.Context) ([]Category, error)
}
