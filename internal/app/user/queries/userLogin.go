package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/pkg/bcrypt"
)

var (
	ErrEmptyLoginCreds  = errors.New("identifier and password must not be empty")
	ErrUserNotFound     = errors.New("user not found")
	ErrPasswordMismatch = errors.New("password does not match")
)

type UserLoginRequest struct {
	Identifier string
	Password   string
}

type UserLoginRequestHandler interface {
	Handle(ctx context.Context, req UserLoginRequest) (*user.User, error)
}

type userLoginRequestHandler struct {
	repo               user.Repository
	encryptionProvider bcrypt.Provider
}

func NewUserLoginHandler(repo user.Repository, encryptionProvider bcrypt.Provider) UserLoginRequestHandler {
	return &userLoginRequestHandler{
		repo:               repo,
		encryptionProvider: encryptionProvider,
	}
}

func (h *userLoginRequestHandler) Handle(ctx context.Context, req UserLoginRequest) (*user.User, error) {
	if req.Identifier == "" || req.Password == "" {
		return nil, ErrEmptyLoginCreds
	}

	user, err := h.repo.GetUserByIdentifier(ctx, req.Identifier)
	if err != nil {
		return nil, ErrUserNotFound
	}

	err = h.encryptionProvider.Matches(user.Password, req.Password)
	if err != nil {
		return nil, ErrPasswordMismatch
	}

	return user, nil
}
