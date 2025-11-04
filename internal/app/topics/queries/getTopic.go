package topicqueries

import (
	"context"

	"github.com/arnald/forum/internal/domain/comment"
	"github.com/arnald/forum/internal/domain/topic"
)

type GetTopicRequest struct {
	UserID  *string `json:"userId"`
	TopicID int     `json:"topicId"`
}

type GetTopicRequestHandler interface {
	Handle(ctx context.Context, req GetTopicRequest) (*topic.Topic, error)
}

type getTopicRequestHandler struct {
	topicRepo   topic.Repository
	commentRepo comment.Repository
}

func NewGetTopicHandler(topicRepo topic.Repository, commentRepo comment.Repository) GetTopicRequestHandler {
	return &getTopicRequestHandler{
		topicRepo:   topicRepo,
		commentRepo: commentRepo,
	}
}

func (h *getTopicRequestHandler) Handle(ctx context.Context, req GetTopicRequest) (*topic.Topic, error) {
	topic, err := h.topicRepo.GetTopicByID(ctx, req.TopicID, req.UserID)
	if err != nil {
		return nil, err
	}

	comments, err := h.commentRepo.GetCommentsWithVotes(ctx, req.TopicID, req.UserID)
	if err != nil {
		return nil, err
	}

	topic.Comments = comments

	return topic, nil
}
