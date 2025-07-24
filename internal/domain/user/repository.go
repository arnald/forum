package user

import (
	"context"
)

type Repository interface {
	GetAll(ctx context.Context) ([]User, error)
	UserRegister(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}
