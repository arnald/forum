package app

import (
	categoryCommands "github.com/arnald/forum/internal/app/categories/commands"
	topicCommands "github.com/arnald/forum/internal/app/topics/commands"
	topicQueries "github.com/arnald/forum/internal/app/topics/queries"
	userQueries "github.com/arnald/forum/internal/app/user/queries"
	"github.com/arnald/forum/internal/domain/categories"
	"github.com/arnald/forum/internal/domain/topic"
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
	CreateTopic    topicCommands.CreateTopicRequestHandler
	UpdateTopic    topicCommands.UpdateTopicRequestHandler
	DeleteTopic    topicCommands.DeleteTopicRequestHandler
	CreateCategory categoryCommands.CreateCategoryRequestHandler
}

type UserServices struct {
	Queries  Queries
	Commands Commands
}

type Services struct {
	UserServices UserServices
}

func NewServices(userRepo user.Repository, categoryRepo categories.Repository, topicRepo topic.Repository) Services {
	uuidProvider := uuid.NewProvider()
	encryption := bcrypt.NewProvider()
	return Services{
		UserServices: UserServices{
			Queries: Queries{
				userQueries.NewUserRegisterHandler(userRepo, uuidProvider, encryption),
				userQueries.NewUserLoginEmailHandler(userRepo, encryption),
				userQueries.NewUserLoginUsernameHandler(userRepo, encryption),
				topicQueries.NewGetTopicHandler(topicRepo),
				topicQueries.NewGetAllTopicsHandler(topicRepo),
			},
			Commands: Commands{
				topicCommands.NewCreateTopicHandler(topicRepo),
				topicCommands.NewUpdateTopicHandler(topicRepo),
				topicCommands.NewDeleteTopicHandler(topicRepo),
				categoryCommands.NewCreateCategoryHandler(categoryRepo),
			},
		},
	}
}
