package updatecomment

import (
	"context"
	"net/http"

	"github.com/arnald/forum/internal/app"
	commentCommands "github.com/arnald/forum/internal/app/comments/commands"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/middleware"
	"github.com/arnald/forum/internal/pkg/helpers"
	"github.com/arnald/forum/internal/pkg/validator"
)

type RequestModel struct {
	CommentID int    `json:"id"`
	Content   string `json:"content"`
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

func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		h.Logger.PrintError(logger.ErrUserNotFoundInContext, nil)
		helpers.RespondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	var commentToUpdate RequestModel

	commentAny, err := helpers.ParseBodyRequest(r, &commentToUpdate)
	if err != nil {
		helpers.RespondWithError(w,
			http.StatusBadRequest,
			"Invalid request payload",
		)

		h.Logger.PrintError(err, nil)
		return
	}
	defer r.Body.Close()

	v := validator.New()

	validator.ValidateUpdateComment(v, commentAny)

	if !v.Valid() {
		helpers.RespondWithError(
			w,
			http.StatusBadRequest,
			v.ToStringErrors(),
		)

		h.Logger.PrintError(logger.ErrValidationFailed, v.Errors)
		return
	}

	_, err = h.UserServices.UserServices.Commands.UpdateComment.Handle(ctx, commentCommands.UpdateCommentRequest{
		CommentID: commentToUpdate.CommentID,
		Content:   commentToUpdate.Content,
		User:      user,
	})
	if err != nil {
		helpers.RespondWithError(w,
			http.StatusInternalServerError,
			"Failed to update comment",
		)

		h.Logger.PrintError(err, nil)
		return
	}

	commentResponse := ResponseModel{
		Message: "Comment updated successfully",
	}

	helpers.RespondWithJSON(
		w,
		http.StatusOK,
		nil,
		commentResponse,
	)

	h.Logger.PrintInfo(
		"Comment updated successfully",
		map[string]string{
			"user_id": user.ID,
		},
	)
}
