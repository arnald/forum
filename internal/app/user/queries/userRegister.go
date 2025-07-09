package queries

import (
	"fmt"
	"time"

	"github.com/arnald/forum/internal/domain/user"
	"github.com/google/uuid"
)

type UserRegisterRequest struct {
	Name     string
	Password string
	Email    string
}

type UserRegisterRequestHandler interface {
	Handle(req UserRegisterRequest) error
}

type userRegisterRequestHandler struct {
	repo user.Repository
}

func NewUserRegisterHandler(repo user.Repository) userRegisterRequestHandler {
	return userRegisterRequestHandler{
		repo: repo,
	}
}

func (h userRegisterRequestHandler) Handle(req UserRegisterRequest) error {
	user := user.User{
		CreatedAt: time.Now(),
		Password:  &req.Password,
		AvatarURL: nil,
		Username:  req.Name,
		Email:     req.Email,
		Role:      "user",
		ID:        uuid.New(),
	}
	err := h.repo.UserRegister(user)
	if err != nil {
		return err
	}
	fmt.Println("User Registered Successfully")

	return nil
}
