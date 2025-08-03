package user

import (
	"context"
)

type Repository interface {
	GetAll(ctx context.Context) ([]User, error)
	UserRegister(ctx context.Context, user *User) error
	CreateSession(session *Session) error
}
