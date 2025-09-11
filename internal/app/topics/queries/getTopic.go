package topicqueries

import (
	"context"

	"github.com/arnald/forum/internal/domain/user"
)

type GetTopicRequest struct {
	TopicID int `json:"topicId"`
}

type GetTopicRequestHandler interface {
	Handle(ctx context.Context, req GetTopicRequest) (*user.Topic, error)
}

type getTopicRequestHandler struct {
	repo user.Repository
}

func NewGetTopicHandler(repo user.Repository) GetTopicRequestHandler {
	return &getTopicRequestHandler{
		repo: repo,
	}
}

func (h *getTopicRequestHandler) Handle(ctx context.Context, req GetTopicRequest) (*user.Topic, error) {
	topic, err := h.repo.GetTopicByID(ctx, req.TopicID)
	if err != nil {
		return nil, err
	}
	return topic, nil
}
