package deletecomment

import (
	"context"
	"net/http"

	"github.com/arnald/forum/internal/app"
	commentCommands "github.com/arnald/forum/internal/app/comments/commands"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/middleware"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type RequestModel struct {
	CommentID int `json:"commentId"`
}

type ResponseModel struct {
	Message string `json:"message"`
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

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	user := middleware.GetUserFromContext(r)

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	var commentToDelete RequestModel

	_, err := helpers.ParseBodyRequest(r, &commentToDelete)
	if err != nil {
		helpers.RespondWithError(w,
			http.StatusBadRequest,
			"Invalid request payload",
		)

		h.Logger.PrintError(err, nil)

		return
	}
	defer r.Body.Close()

	err = h.UserServices.UserServices.Commands.DeleteComment.Handle(ctx, commentCommands.DeleteCommentRequest{
		CommentID: commentToDelete.CommentID,
		User:      user,
	})
	if err != nil {
		helpers.RespondWithError(w,
			http.StatusInternalServerError,
			"Failed to delete comment",
		)

		h.Logger.PrintError(err, nil)
		return
	}

	commentResponse := ResponseModel{
		Message: "Comment deleted successfully",
	}

	helpers.RespondWithJSON(
		w,
		http.StatusOK,
		nil,
		commentResponse,
	)

	h.Logger.PrintInfo(
		"Comment deleted successfully",
		map[string]string{
			"user_id": user.ID,
		},
	)
}
