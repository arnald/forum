//nolint:dupl
package userqueries

import (
	"context"

	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/pkg/bcrypt"
)

type UserLoginUsernameRequest struct {
	Username string
	Password string
}

type UserLoginUsernameRequestHandler interface {
	Handle(ctx context.Context, req UserLoginUsernameRequest) (*user.User, error)
}

type userLoginUsernameRequestHandler struct {
	repo               user.Repository
	encryptionProvider bcrypt.Provider
}

func NewUserLoginUsernameHandler(repo user.Repository, encryptionProvider bcrypt.Provider) UserLoginUsernameRequestHandler {
	return &userLoginUsernameRequestHandler{
		repo:               repo,
		encryptionProvider: encryptionProvider,
	}
}

func (h *userLoginUsernameRequestHandler) Handle(ctx context.Context, req UserLoginUsernameRequest) (*user.User, error) {
	user, err := h.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	err = h.encryptionProvider.Matches(user.Password, req.Password)
	if err != nil {
		return nil, ErrPasswordMismatch
	}

	return user, nil
}
