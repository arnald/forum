package user

import (
	"context"
)

type Repository interface {
	GetAll(ctx context.Context) ([]User, error)
	UserRegister(user User) error
}
