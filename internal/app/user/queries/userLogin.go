package userqueries

import (
	"context"
	"fmt"

	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/pkg/bcrypt"
)

type UserLoginRequest struct {
	Email    string
	Password string
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
	user, err := h.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	err = h.encryptionProvider.Matches(user.Password, req.Password)
	if err != nil {
		return nil, fmt.Errorf("password does not match: %w", err)
	}

	return user, nil
}
