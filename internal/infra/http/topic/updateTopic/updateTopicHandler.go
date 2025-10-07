package updatetopic

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

type RequestModel struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	ImagePath  string `json:"imagePath"`
	CategoryID int    `json:"categoryId"`
	TopicID    int    `json:"topicId"`
}

type ResponseModel struct {
	UserID  string `json:"userId"`
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

func (h *Handler) UpdateTopic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	user := middleware.GetUserFromContext(r)

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	var topicToUpdate RequestModel

	topicAny, err := helpers.ParseBodyRequest(r, &topicToUpdate)
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

	validator.ValidateCreateTopic(v, topicAny)

	if !v.Valid() {
		helpers.RespondWithError(
			w,
			http.StatusBadRequest,
			v.ToStringErrors(),
		)

		h.Logger.PrintError(logger.ErrValidationFailed, v.Errors)
		return
	}

	topic, err := h.UserServices.UserServices.Commands.UpdateTopic.Handle(ctx, topicCommands.UpdateTopicRequest{
		CategoryID: topicToUpdate.CategoryID,
		TopicID:    topicToUpdate.TopicID,
		Title:      topicToUpdate.Title,
		Content:    topicToUpdate.Content,
		ImagePath:  topicToUpdate.ImagePath,
		User:       user,
	})
	if err != nil {
		helpers.RespondWithError(w,
			http.StatusInternalServerError,
			"Failed to create topic",
		)

		h.Logger.PrintError(err, nil)

		return
	}

	topicResponse := ResponseModel{
		UserID:  topic.UserID,
		Message: "Topic updated successfully",
	}

	helpers.RespondWithJSON(
		w,
		http.StatusCreated,
		nil,
		topicResponse,
	)

	h.Logger.PrintInfo(
		"Topic updated successfully",
		map[string]string{
			"user_id": topicResponse.UserID,
		},
	)
}
