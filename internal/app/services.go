package app

import (
	topicCommands "github.com/arnald/forum/internal/app/topics/commands"
	topicQueries "github.com/arnald/forum/internal/app/topics/queries"
	userQueries "github.com/arnald/forum/internal/app/user/queries"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/pkg/bcrypt"
	"github.com/arnald/forum/internal/pkg/uuid"
)

type Queries struct {
	UserRegister      userQueries.UserRegisterRequestHandler
	UserLoginEmail    userQueries.UserLoginEmailRequestHandler
	UserLoginUsername userQueries.UserLoginUsernameRequestHandler
	GetTopic          topicQueries.GetTopicRequestHandler
	GetAllTopics      topicQueries.GetAllTopicsRequestHandler
}

type Commands struct {
	CreateTopic topicCommands.CreateTopicRequestHandler
	UpdateTopic topicCommands.UpdateTopicRequestHandler
	DeleteTopic topicCommands.DeleteTopicRequestHandler
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
				userQueries.NewUserRegisterHandler(repo, uuidProvider, encryption),
				userQueries.NewUserLoginEmailHandler(repo, encryption),
				userQueries.NewUserLoginUsernameHandler(repo, encryption),
				topicQueries.NewGetTopicHandler(repo),
				topicQueries.NewGetAllTopicsHandler(repo),
			},
			Commands: Commands{
				topicCommands.NewCreateTopicHandler(repo),
				topicCommands.NewUpdateTopicHandler(repo),
				topicCommands.NewDeleteTopicHandler(repo),
			},
		},
	}
}
