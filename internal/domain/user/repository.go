package user

import (
	"context"
)

type Repository interface {
	GetAll(ctx context.Context) ([]User, error)
}
