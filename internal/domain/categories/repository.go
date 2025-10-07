package categories

import "context"

type Repository interface {
	CreateCategory(ctx context.Context, category *Category) error
}
