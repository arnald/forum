package castvote

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/arnald/forum/internal/app"
	commentqueries "github.com/arnald/forum/internal/app/comments/queries"
	topicqueries "github.com/arnald/forum/internal/app/topics/queries"
	votecommands "github.com/arnald/forum/internal/app/votes/commands"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/comment"
	"github.com/arnald/forum/internal/domain/notification"
	"github.com/arnald/forum/internal/domain/topic"
	"github.com/arnald/forum/internal/domain/vote"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/middleware"
	"github.com/arnald/forum/internal/infra/storage/notifications"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type RequestModel struct {
	TopicID      *int `json:"topicId"`
	CommentID    *int `json:"commentId,omitempty"`
	ReactionType int  `json:"reactionType"`
}

type ResponseModel struct {
	Counts  *vote.Counts `json:"counts"`
	Message string       `json:"message"`
}

type Handler struct {
	Services      app.Services
	Config        *config.ServerConfig
	Logger        logger.Logger
	Notifications *notifications.NotificationService
}

func NewHandler(services app.Services, config *config.ServerConfig, logger logger.Logger, notifications *notifications.NotificationService) *Handler {
	return &Handler{
		Services:      services,
		Config:        config,
		Logger:        logger,
		Notifications: notifications,
	}
}

func (h *Handler) CastVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(
			w,
			http.StatusMethodNotAllowed,
			"Invalid request method",
		)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		h.Logger.PrintError(logger.ErrUserNotFoundInContext, nil)
		helpers.RespondWithError(
			w,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	var req RequestModel
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(
			w,
			http.StatusBadRequest,
			"invalid JSON",
		)
		return
	}

	target := vote.Target{
		TopicID:   req.TopicID,
		CommentID: req.CommentID,
	}

	err = h.Services.UserServices.Commands.CastVote.Handle(ctx, votecommands.CastVoteRequest{
		UserID:       user.ID,
		Target:       target,
		ReactionType: req.ReactionType,
	})
	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(
			w,
			http.StatusInternalServerError,
			"failed to cast vote",
		)
		return
	}

	h.sendVoteNotification(
		ctx,
		user.Username,
		user.ID,
		req,
	)

	Response := ResponseModel{
		Message: "Vote cast successfully",
	}

	helpers.RespondWithJSON(w,
		http.StatusOK,
		nil,
		Response,
	)
}

func (h *Handler) getCommentOwner(ctx context.Context, commentID int) (*comment.Comment, error) {
	return h.Services.UserServices.Queries.GetComment.Handle(ctx, commentqueries.GetCommentRequest{
		CommentID: commentID,
	})
}

func (h *Handler) getTopicOwner(ctx context.Context, topicID int) (*topic.Topic, error) {
	return h.Services.UserServices.Queries.GetTopic.Handle(ctx, topicqueries.GetTopicRequest{
		TopicID: topicID,
	})
}

func (h *Handler) sendVoteNotification(ctx context.Context, username, userID string, req RequestModel) {
	var ownerID string
	var contentType string
	var contentID string

	if req.CommentID != nil {
		comment, err := h.getCommentOwner(ctx, *req.CommentID)
		if err != nil {
			h.Logger.PrintError(err, nil)
			return
		}
		ownerID = comment.UserID
		contentType = "comment"
		contentID = strconv.Itoa(*req.CommentID)
	} else if req.TopicID != nil {
		topic, err := h.getTopicOwner(ctx, *req.TopicID)
		if err != nil {
			h.Logger.PrintError(err, nil)
			return
		}
		ownerID = topic.UserID
		contentType = "topic"
		contentID = strconv.Itoa(*req.TopicID)
	}

	if ownerID == "" || ownerID == userID {
		return
	}

	var message string
	var title string
	var notificationType notification.Type
	if req.ReactionType == 1 {
		message = fmt.Sprintf("%s liked your %s", username, contentType)
		title = "New like!"
		notificationType = notification.NotificationTypeLike
	} else if req.ReactionType == -1 {
		message = fmt.Sprintf("%s disliked your %s", username, contentType)
		title = "New dislike!"
		notificationType = notification.NotificationTypeDislike
	}

	notification := &notification.Notification{
		Type:        notificationType,
		RelatedID:   contentID,
		UserID:      ownerID,
		ActorID:     userID,
		Message:     message,
		Title:       title,
		RelatedType: contentType,
	}

	err := h.Notifications.CreateNotification(ctx, notification)
	if err != nil {
		h.Logger.PrintError(err, nil)
	}
}
