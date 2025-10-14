package topicqueries

import (
	"context"

	"github.com/arnald/forum/internal/domain/topic"
)

type GetTopicRequest struct {
	TopicID int `json:"topicId"`
}

type GetTopicRequestHandler interface {
	Handle(ctx context.Context, req GetTopicRequest) (*topic.Topic, error)
}

type getTopicRequestHandler struct {
	repo topic.Repository
}

func NewGetTopicHandler(repo topic.Repository) GetTopicRequestHandler {
	return &getTopicRequestHandler{
		repo: repo,
	}
}

func (h *getTopicRequestHandler) Handle(ctx context.Context, req GetTopicRequest) (*topic.Topic, error) {
	topic, err := h.repo.GetTopicByID(ctx, req.TopicID)
	if err != nil {
		return nil, err
	}
	return topic, nil
}
