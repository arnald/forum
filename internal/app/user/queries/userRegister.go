package queries

import (
	"fmt"
	"time"

	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/pkg/bcrypt"
	"github.com/arnald/forum/internal/pkg/uuid"
)

type UserRegisterRequest struct {
	Name     string
	Password string
	Email    string
}

type UserRegisterRequestHandler interface {
	Handle(req UserRegisterRequest) (*user.User, error)
}

type userRegisterRequestHandler struct {
	uuidiProvider      uuid.Provider
	encryptionProvider bcrypt.Provider
	repo               user.Repository
}

func NewUserRegisterHandler(repo user.Repository, uuidProvider uuid.Provider, en bcrypt.Provider) userRegisterRequestHandler {
	return userRegisterRequestHandler{
		repo:               repo,
		uuidiProvider:      uuidProvider,
		encryptionProvider: en,
	}
}

func (h userRegisterRequestHandler) Handle(req UserRegisterRequest) (*user.User, error) {
	user := &user.User{
		CreatedAt: time.Now(),
		Password:  &req.Password,
		AvatarURL: nil,
		Username:  req.Name,
		Email:     req.Email,
		ID:        h.uuidiProvider.NewUUID(),
	}

	encryptedPass, err := h.encryptionProvider.Generate(*user.Password)
	if err != nil {
		return nil, err
	}

	err = h.repo.UserRegister(user, encryptedPass)
	if err != nil {
		return nil, err
	}
	fmt.Println("User Registered Successfully")

	return user, err
}
