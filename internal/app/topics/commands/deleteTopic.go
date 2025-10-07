package topiccommands

import (
	"context"

	"github.com/arnald/forum/internal/domain/user"
)

type DeleteTopicRequest struct {
	User    *user.User
	TopicID int `json:"topicId"`
}

type DeleteTopicRequestHandler interface {
	Handle(ctx context.Context, req DeleteTopicRequest) error
}

type deleteTopicRequestHandler struct {
	repo user.Repository
}

func NewDeleteTopicHandler(repo user.Repository) DeleteTopicRequestHandler {
	return &deleteTopicRequestHandler{
		repo: repo,
	}
}

func (h *deleteTopicRequestHandler) Handle(ctx context.Context, req DeleteTopicRequest) error {
	err := h.repo.DeleteTopic(ctx, req.User.ID, req.TopicID)
	if err != nil {
		return err
	}
	return nil
}
