package gettopic

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/arnald/forum/internal/app"
	topicQueries "github.com/arnald/forum/internal/app/topics/queries"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/comment"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type RequestModel struct {
	TopicID int `json:"topicId"`
}

type ResponseModel struct {
	Title      string            `json:"title"`
	Content    string            `json:"content"`
	ImagePath  string            `json:"imagePath"`
	UserID     string            `json:"userId"`
	CreatedAt  string            `json:"createdAt"`
	UpdatedAt  string            `json:"updatedAt"`
	Comments   []comment.Comment `json:"comments,omitempty"`
	TopicID    int               `json:"topicId"`
	CategoryID int               `json:"categoryId"`
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

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	var topicToGet RequestModel

	err := json.NewDecoder(r.Body).Decode(&topicToGet)
	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	topic, err := h.UserServices.UserServices.Queries.GetTopic.Handle(ctx, topicQueries.GetTopicRequest{
		TopicID: topicToGet.TopicID,
	})
	if err != nil {
		if errors.Is(err, topicQueries.ErrTopicNotFound) {
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
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil, response)
}
