package createcomment

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
	TopicID int    `json:"topicId"`
	Content string `json:"content"`
}

type ResponseModel struct {
	CommentID int    `json:"commentId"`
	Message   string `json:"message"`
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

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	user := middleware.GetUserFromContext(r)

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	var commentToCreate RequestModel

	commentAny, err := helpers.ParseBodyRequest(r, &commentToCreate)
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

	validator.ValidateCreateComment(v, commentAny)

	if !v.Valid() {
		helpers.RespondWithError(
			w,
			http.StatusBadRequest,
			v.ToStringErrors(),
		)

		h.Logger.PrintError(logger.ErrValidationFailed, v.Errors)
		return
	}

	comment, err := h.UserServices.UserServices.Commands.CreateComment.Handle(ctx, commentCommands.CreateCommentRequest{
		TopicID: commentToCreate.TopicID,
		Content: commentToCreate.Content,
		User:    user,
	})
	if err != nil {
		helpers.RespondWithError(w,
			http.StatusInternalServerError,
			"Failed to create comment",
		)

		h.Logger.PrintError(err, nil)
		return
	}

	commentResponse := ResponseModel{
		CommentID: comment.ID,
		Message:   "Comment created successfully",
	}

	helpers.RespondWithJSON(
		w,
		http.StatusCreated,
		nil,
		commentResponse,
	)

	h.Logger.PrintInfo(
		"Comment created successfully",
		map[string]string{
			"user_id":    user.ID,
			"comment_id": string(comment.ID),
		},
	)
}
