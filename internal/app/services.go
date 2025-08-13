package app

import (
	"github.com/arnald/forum/internal/app/topics/commands"
	"github.com/arnald/forum/internal/app/user/queries"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/pkg/bcrypt"
	"github.com/arnald/forum/internal/pkg/uuid"
)

type Queries struct {
	UserRegister      queries.UserRegisterRequestHandler
	UserLoginEmail    queries.UserLoginEmailRequestHandler
	UserLoginUsername queries.UserLoginUsernameRequestHandler
}

type Commands struct {
	CreateTopic commands.CreateTopicRequestHandler
	UpdateTopic commands.UpdateTopicRequestHandler
}

type UserServices struct {
	Queries  Queries
	Commands Commands
}

type Services struct {
	UserServices UserServices
}

func NewServices(repo user.Repository) Services {
	uuidProvider := uuid.NewProvider()
	encryption := bcrypt.NewProvider()
	return Services{
		UserServices: UserServices{
			Queries: Queries{
				queries.NewUserRegisterHandler(repo, uuidProvider, encryption),
				queries.NewUserLoginEmailHandler(repo, encryption),
				queries.NewUserLoginUsernameHandler(repo, encryption),
			},
			Commands: Commands{
				commands.NewCreateTopicHandler(repo),
				commands.NewUpdateTopicHandler(repo),
			},
		},
	}
}
