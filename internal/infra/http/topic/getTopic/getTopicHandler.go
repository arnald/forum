package gettopic

import (
	"context"
	"errors"
	"net/http"

	"github.com/arnald/forum/internal/app"
	topicQueries "github.com/arnald/forum/internal/app/topics/queries"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/comment"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/middleware"
	"github.com/arnald/forum/internal/infra/storage/sqlite/topics"
	"github.com/arnald/forum/internal/pkg/helpers"
	"github.com/arnald/forum/internal/pkg/validator"
)

type ResponseModel struct {
	Title      string            `json:"title"`
	Content    string            `json:"content"`
	ImagePath  string            `json:"imagePath"`
	UserID     string            `json:"userId"`
	CreatedAt  string            `json:"createdAt"`
	UpdatedAt  string            `json:"updatedAt"`
	Comments   []comment.Comment `json:"comments"`
	TopicID    int               `json:"topicId"`
	CategoryID int               `json:"categoryId"`
	Upvotes    int               `json:"upvotes"`
	Downvotes  int               `json:"downvotes"`
	Score      int               `json:"score"`
	UserVote   *int              `json:"userVote"`
}

type Handler struct {
	UserServices app.Services
	Config       *config.ServerConfig
	Logger       logger.Logger
}

func NewHandler(userServices app.Services, config *config.ServerConfig, logger logger.Logger) *Handler {
	return &Handler{
		UserServices: userServices,
		Config:       config,
		Logger:       logger,
	}
}

func (h *Handler) GetTopic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	var userID *string
	user := middleware.GetUserFromContext(r)
	if user != nil {
		userID = &user.ID
	}

	topicID, err := helpers.GetQueryInt(r, "id")
	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	val := validator.New()

	testStruct := &struct {
		TopicID int
	}{
		TopicID: topicID,
	}
	validator.ValidateGetTopic(val, testStruct)

	if !val.Valid() {
		h.Logger.PrintError(logger.ErrValidationFailed, val.Errors)
		helpers.RespondWithError(w,
			http.StatusBadRequest,
			val.ToStringErrors(),
		)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	topic, err := h.UserServices.UserServices.Queries.GetTopic.Handle(ctx, topicQueries.GetTopicRequest{
		TopicID: topicID,
		UserID:  userID,
	})
	if err != nil {
		if errors.Is(err, topics.ErrTopicNotFound) {
			helpers.RespondWithError(w, http.StatusNotFound, "Topic not found")
			return
		}

		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		h.Logger.PrintError(err, nil)
		return
	}

	response := ResponseModel{
		TopicID:    topic.ID,
		CategoryID: topic.CategoryID,
		Title:      topic.Title,
		Content:    topic.Content,
		ImagePath:  topic.ImagePath,
		UserID:     topic.UserID,
		CreatedAt:  topic.CreatedAt,
		UpdatedAt:  topic.UpdatedAt,
		Comments:   topic.Comments,
		Upvotes:    topic.UpvoteCount,
		Downvotes:  topic.DownvoteCount,
		Score:      topic.VoteScore,
		UserVote:   topic.UserVote,
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil, response)
}
