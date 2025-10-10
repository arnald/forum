package getcomment

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/arnald/forum/internal/app"
	commentQueries "github.com/arnald/forum/internal/app/comments/queries"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type RequestModel struct {
	CommentID int `json:"commentId"`
}

type ResponseModel struct {
	ID        int    `json:"id"`
	UserID    string `json:"userId"`
	Username  string `json:"username"`
	TopicID   int    `json:"topicId"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
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

func (h *Handler) GetComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	var commentToGet RequestModel

	err := json.NewDecoder(r.Body).Decode(&commentToGet)
	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	comment, err := h.UserServices.UserServices.Queries.GetComment.Handle(ctx, commentQueries.GetCommentRequest{
		CommentID: commentToGet.CommentID,
	})
	if err != nil {
		if errors.Is(err, commentQueries.ErrCommentNotFound) {
			helpers.RespondWithError(w, http.StatusNotFound, "Comment not found")
			return
		}

		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		h.Logger.PrintError(err, nil)
		return
	}

	response := ResponseModel{
		ID:        comment.ID,
		UserID:    comment.UserID,
		Username:  comment.OwnerUsername,
		TopicID:   comment.TopicID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: comment.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil, response)
}
