//nolint:dupl
package queries

import (
	"context"
	"errors"

	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/pkg/bcrypt"
)

var ErrEmptyEmail = errors.New("email must not be empty")

type UserLoginEmailRequest struct {
	Email    string
	Password string
}

type UserLoginEmailRequestHandler interface {
	Handle(ctx context.Context, req UserLoginEmailRequest) (*user.User, error)
}

type userLoginEmailRequestHandler struct {
	repo               user.Repository
	encryptionProvider bcrypt.Provider
}

func NewUserLoginEmailHandler(repo user.Repository, encryptionProvider bcrypt.Provider) UserLoginEmailRequestHandler {
	return &userLoginEmailRequestHandler{
		repo:               repo,
		encryptionProvider: encryptionProvider,
	}
}

func (h *userLoginEmailRequestHandler) Handle(ctx context.Context, req UserLoginEmailRequest) (*user.User, error) {
	if req.Email == "" || req.Password == "" {
		return nil, ErrEmptyLoginCreds
	}

	user, err := h.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	err = h.encryptionProvider.Matches(user.Password, req.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}
