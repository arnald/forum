package app

import (
	categoryCommands "github.com/arnald/forum/internal/app/categories/commands"
	categoryQueries "github.com/arnald/forum/internal/app/categories/queries"
	topicCommands "github.com/arnald/forum/internal/app/topics/commands"
	topicQueries "github.com/arnald/forum/internal/app/topics/queries"
	userCommands "github.com/arnald/forum/internal/app/user/commands"
	userQueries "github.com/arnald/forum/internal/app/user/queries"
	"github.com/arnald/forum/internal/domain/category"
	"github.com/arnald/forum/internal/domain/topic"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/pkg/bcrypt"
	"github.com/arnald/forum/internal/pkg/uuid"
)

type Queries struct {
	GetTopic          topicQueries.GetTopicRequestHandler
	GetAllTopics      topicQueries.GetAllTopicsRequestHandler
	UserLoginEmail    userQueries.UserLoginEmailRequestHandler
	UserLoginUsername userQueries.UserLoginUsernameRequestHandler
	GetCategoryByID   categoryQueries.GetCategoryByIDHandler
}

type Commands struct {
	UserRegister   userCommands.UserRegisterRequestHandler
	CreateTopic    topicCommands.CreateTopicRequestHandler
	UpdateTopic    topicCommands.UpdateTopicRequestHandler
	DeleteTopic    topicCommands.DeleteTopicRequestHandler
	CreateCategory categoryCommands.CreateCategoryRequestHandler
	DeleteCategory categoryCommands.DeleteCategoryRequestHandler
	UpdateCategory categoryCommands.UpdateCategoryRequestHandler
}

type UserServices struct {
	Queries  Queries
	Commands Commands
}

type Services struct {
	UserServices UserServices
}

func NewServices(userRepo user.Repository, categoryRepo category.Repository, topicRepo topic.Repository) Services {
	uuidProvider := uuid.NewProvider()
	encryption := bcrypt.NewProvider()
	return Services{
		UserServices: UserServices{
			Queries: Queries{
				topicQueries.NewGetTopicHandler(topicRepo),
				topicQueries.NewGetAllTopicsHandler(topicRepo),
				userQueries.NewUserLoginEmailHandler(userRepo, encryption),
				userQueries.NewUserLoginUsernameHandler(userRepo, encryption),
				categoryQueries.NewGetCategoryByIDHandler(categoryRepo),
			},
			Commands: Commands{
				userCommands.NewUserRegisterHandler(userRepo, uuidProvider, encryption),
				topicCommands.NewCreateTopicHandler(topicRepo),
				topicCommands.NewUpdateTopicHandler(topicRepo),
				topicCommands.NewDeleteTopicHandler(topicRepo),
				categoryCommands.NewCreateCategoryHandler(categoryRepo),
				categoryCommands.NewDeleteCategoryHandler(categoryRepo),
				categoryCommands.NewUpdateCategoryHandler(categoryRepo),
			},
		},
	}
}
