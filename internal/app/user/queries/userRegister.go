package userqueries

import (
	"context"
	"time"

	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/pkg/bcrypt"
	"github.com/arnald/forum/internal/pkg/helpers"
	"github.com/arnald/forum/internal/pkg/uuid"
)

type UserRegisterRequest struct {
	Name     string
	Password string
	Email    string
}

type UserRegisterRequestHandler interface {
	Handle(ctx context.Context, req UserRegisterRequest) (*user.User, error)
}

type userRegisterRequestHandler struct {
	uuidiProvider      uuid.Provider
	encryptionProvider bcrypt.Provider
	repo               user.Repository
}

func NewUserRegisterHandler(repo user.Repository, uuidProvider uuid.Provider, en bcrypt.Provider) UserRegisterRequestHandler {
	return userRegisterRequestHandler{
		repo:               repo,
		uuidiProvider:      uuidProvider,
		encryptionProvider: en,
	}
}

func (h userRegisterRequestHandler) Handle(ctx context.Context, req UserRegisterRequest) (*user.User, error) {
	user := &user.User{
		CreatedAt: time.Now(),
		Password:  req.Password,
		AvatarURL: nil,
		Username:  req.Name,
		Email:     req.Email,
		ID:        h.uuidiProvider.NewUUID(),
	}

	err := helpers.ValidateEmail(user.Email)
	if err != nil {
		return nil, err
	}

	encryptedPass, err := h.encryptionProvider.Generate(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = encryptedPass

	err = h.repo.UserRegister(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, err
}
