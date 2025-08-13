package deletetopic

import (
	"context"
	"net/http"

	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/app/topics/commands"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/middleware"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type DeleteTopicRequestModel struct {
	TopicID int `json:"topic_id"`
}

type DeleteTopicResponseModel struct {
	UserID  string `json:"user_id"`
	TopicID int    `json:"topic_id"`
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

func (h *Handler) DeleteTopic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	user := middleware.GetUserFromContext(r)

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	var topicToDelete DeleteTopicRequestModel

	_, err := helpers.ParseBodyRequest(r, &topicToDelete)
	if err != nil {
		helpers.RespondWithError(w,
			http.StatusBadRequest,
			"Invalid request payload",
		)

		h.Logger.PrintError(err, nil)

		return
	}
	defer r.Body.Close()

	err = h.UserServices.UserServices.Commands.DeleteTopic.Handle(ctx, commands.DeleteTopicRequest{
		TopicID: topicToDelete.TopicID,
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
	if err != nil {
		helpers.RespondWithError(w,
			http.StatusInternalServerError,
			"Failed to delete topic",
		)

		h.Logger.PrintError(err, nil)

		return
	}

	topicResponse := DeleteTopicResponseModel{
		UserID:  user.ID,
		TopicID: topicToDelete.TopicID,
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
