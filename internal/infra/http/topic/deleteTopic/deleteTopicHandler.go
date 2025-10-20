package deletetopic

import (
	"context"
	"net/http"

	"github.com/arnald/forum/internal/app"
	topicCommands "github.com/arnald/forum/internal/app/topics/commands"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/middleware"
	"github.com/arnald/forum/internal/pkg/helpers"
	"github.com/arnald/forum/internal/pkg/validator"
)

type ResponseModel struct {
	UserID  string `json:"userId"`
	Message string `json:"message"`
	TopicID int    `json:"topicId"`
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

func (h *Handler) DeleteTopic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	user := middleware.GetUserFromContext(r)

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

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

	validator.ValidateDeleteTopic(val, &struct {
		TopicID int
	}{
		TopicID: topicID,
	})

	if !val.Valid() {
		h.Logger.PrintError(logger.ErrValidationFailed, val.Errors)
		helpers.RespondWithError(w,
			http.StatusBadRequest,
			val.ToStringErrors())
		return
	}

	err = h.UserServices.UserServices.Commands.DeleteTopic.Handle(ctx, topicCommands.DeleteTopicRequest{
		TopicID: topicID,
		User:    user,
	})
	if err != nil {
		helpers.RespondWithError(w,
			http.StatusInternalServerError,
			"Failed to delete topic",
		)

		h.Logger.PrintError(err, nil)

		return
	}

	topicResponse := ResponseModel{
		UserID:  user.ID,
		TopicID: topicID,
		Message: "Topic deleted successfully",
	}

	helpers.RespondWithJSON(
		w,
		http.StatusOK,
		nil,
		topicResponse,
	)

	h.Logger.PrintInfo(
		"Topic deleted successfully",
		map[string]string{
			"user_id": topicResponse.UserID,
		},
	)
}
