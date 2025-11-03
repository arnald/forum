package createcomment

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/arnald/forum/internal/app"
	commentCommands "github.com/arnald/forum/internal/app/comments/commands"
	topicqueries "github.com/arnald/forum/internal/app/topics/queries"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/notification"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/middleware"
	"github.com/arnald/forum/internal/infra/storage/notifications"
	"github.com/arnald/forum/internal/pkg/helpers"
	"github.com/arnald/forum/internal/pkg/validator"
)

type RequestModel struct {
	Content string `json:"content"`
	TopicID int    `json:"topicId"`
}

type ResponseModel struct {
	Message   string `json:"message"`
	CommentID int    `json:"commentId"`
}

type Handler struct {
	UserServices app.Services
	Config       *config.ServerConfig
	Logger       logger.Logger
	Notification *notifications.NotificationService
}

func NewHandler(userServices app.Services, config *config.ServerConfig, logger logger.Logger, notifications *notifications.NotificationService) *Handler {
	return &Handler{
		UserServices: userServices,
		Config:       config,
		Logger:       logger,
		Notification: notifications,
	}
}

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	topic, err := h.UserServices.UserServices.Queries.GetTopic.Handle(ctx, topicqueries.GetTopicRequest{
		UserID:  &user.ID,
		TopicID: comment.TopicID,
	})
	if err != nil {
		h.Logger.PrintError(err, nil)
	}

	notification := &notification.Notification{
		ActorID:     user.Username,
		UserID:      topic.UserID,
		RelatedID:   strconv.Itoa(comment.TopicID),
		RelatedType: "topic",
		Type:        notification.NotificationTypeReply,
		Title:       "New comment",
		Message:     fmt.Sprintf("%s commented on your Topic %s", user.Username, topic.Title),
	}

	err = h.Notification.CreateNotification(ctx, notification)
	if err != nil {
		h.Logger.PrintError(err, nil)
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
			"comment_id": strconv.Itoa(comment.ID),
		},
	)
}
