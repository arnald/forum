package app

import (
	categoryCommands "github.com/arnald/forum/internal/app/categories/commands"
	categoryQueries "github.com/arnald/forum/internal/app/categories/queries"
	commentCommands "github.com/arnald/forum/internal/app/comments/commands"
	commentQueries "github.com/arnald/forum/internal/app/comments/queries"
	topicCommands "github.com/arnald/forum/internal/app/topics/commands"
	topicQueries "github.com/arnald/forum/internal/app/topics/queries"
	userCommands "github.com/arnald/forum/internal/app/user/commands"
	userQueries "github.com/arnald/forum/internal/app/user/queries"
	"github.com/arnald/forum/internal/domain/category"
	"github.com/arnald/forum/internal/domain/comment"
	"github.com/arnald/forum/internal/domain/topic"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/pkg/bcrypt"
	"github.com/arnald/forum/internal/pkg/uuid"
)

type Queries struct {
	GetTopic           topicQueries.GetTopicRequestHandler
	GetAllTopics       topicQueries.GetAllTopicsRequestHandler
	GetComment         commentQueries.GetCommentRequestHandler
	GetCommentsByTopic commentQueries.GetCommentsByTopicRequestHandler
	UserLoginEmail     userQueries.UserLoginEmailRequestHandler
	UserLoginUsername  userQueries.UserLoginUsernameRequestHandler
	GetCategoryByID    categoryQueries.GetCategoryByIDHandler
	GetAllCategories   categoryQueries.GetAllCategoriesRequestHandler
}

type Commands struct {
	UserRegister   userCommands.UserRegisterRequestHandler
	CreateTopic    topicCommands.CreateTopicRequestHandler
	UpdateTopic    topicCommands.UpdateTopicRequestHandler
	DeleteTopic    topicCommands.DeleteTopicRequestHandler
	CreateComment  commentCommands.CreateCommentRequestHandler
	UpdateComment  commentCommands.UpdateCommentRequestHandler
	DeleteComment  commentCommands.DeleteCommentRequestHandler
	CreateCategory categoryCommands.CreateCategoryRequestHandler
	UpdateCategory categoryCommands.UpdateCategoryRequestHandler
	DeleteCategory categoryCommands.DeleteCategoryRequestHandler
}

type UserServices struct {
	Queries  Queries
	Commands Commands
}

type Services struct {
	UserServices UserServices
}

func NewServices(userRepo user.Repository, categoryRepo category.Repository, topicRepo topic.Repository, commentRepo comment.Repository) Services {
	uuidProvider := uuid.NewProvider()
	encryption := bcrypt.NewProvider()
	return Services{
		UserServices: UserServices{
			Queries: Queries{
				topicQueries.NewGetTopicHandler(topicRepo),
				topicQueries.NewGetAllTopicsHandler(topicRepo),
				commentQueries.NewGetCommentHandler(commentRepo),
				commentQueries.NewGetCommentsByTopicRequestHandler(commentRepo),
				userQueries.NewUserLoginEmailHandler(userRepo, encryption),
				userQueries.NewUserLoginUsernameHandler(userRepo, encryption),
				categoryQueries.NewGetCategoryByIDHandler(categoryRepo),
				categoryQueries.NewGetAllCategoriesHandler(categoryRepo),
			},
			Commands: Commands{
				userCommands.NewUserRegisterHandler(userRepo, uuidProvider, encryption),
				topicCommands.NewCreateTopicHandler(topicRepo),
				topicCommands.NewUpdateTopicHandler(topicRepo),
				topicCommands.NewDeleteTopicHandler(topicRepo),
				commentCommands.NewCreateCommentRequestHandler(commentRepo),
				commentCommands.NewUpdateCommentRequestHandler(commentRepo),
				commentCommands.NewDeleteCommentHandler(commentRepo),
				categoryCommands.NewCreateCategoryHandler(categoryRepo),
				categoryCommands.NewUpdateCategoryHandler(categoryRepo),
				categoryCommands.NewDeleteCategoryHandler(categoryRepo),
			},
		},
	}
}
