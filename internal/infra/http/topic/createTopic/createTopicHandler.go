package createtopic

import (
	"context"
	"net/http"

	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/app/topics/commands"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/middleware"
	"github.com/arnald/forum/internal/pkg/helpers"
	"github.com/arnald/forum/internal/pkg/validator"
)

type CreateTopicRequestModel struct {
	CategoryID int    `json:"category_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	ImagePath  string `json:"image_path"`
}

type CreateTopicResponseModel struct {
	UserID  string `json:"user_id"`
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

func (h *Handler) CreateTopic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	user := middleware.GetUserFromContext(r)

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	var topicToCreate CreateTopicRequestModel

	topicAny, err := helpers.ParseBodyRequest(r, &topicToCreate)
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

	topic, err := h.UserServices.UserServices.Commands.CreateTopic.Handle(ctx, commands.CreateTopicRequest{
		CategoryID: topicToCreate.CategoryID,
		Title:      topicToCreate.Title,
		Content:    topicToCreate.Content,
		ImagePath:  topicToCreate.ImagePath,
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

	topicResponse := CreateTopicResponseModel{
		UserID:  topic.UserID,
		Message: "Topic created successfully",
	}

	helpers.RespondWithJSON(
		w,
		http.StatusCreated,
		nil,
		topicResponse,
	)

	h.Logger.PrintInfo(
		"Topic created successfully",
		map[string]string{
			"user_id": topicResponse.UserID,
		},
	)
}
