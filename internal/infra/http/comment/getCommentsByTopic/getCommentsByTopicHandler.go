package getcommentsbytopic

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/arnald/forum/internal/app"
	commentQueries "github.com/arnald/forum/internal/app/comments/queries"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/comment"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type RequestModel struct {
	TopicID int `json:"topicId"`
}

type ResponseModel struct {
	Comments []comment.Comment `json:"comments"`
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

func (h *Handler) GetCommentsByTopic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	var req RequestModel

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	comments, err := h.UserServices.UserServices.Queries.GetCommentsByTopic.Handle(ctx, commentQueries.GetCommentsByTopicRequest{
		TopicID: req.TopicID,
	})
	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to get comments")
		return
	}

	response := ResponseModel{
		Comments: comments,
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil, response)
}
