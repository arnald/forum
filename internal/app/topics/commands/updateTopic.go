package topiccommands

import (
	"context"

	"github.com/arnald/forum/internal/domain/user"
)

type UpdateTopicRequest struct {
	User       *user.User
	Title      string `json:"title"`
	Content    string `json:"content"`
	ImagePath  string `json:"imagePath"`
	TopicID    int    `json:"topicId"`
	CategoryID int    `json:"categoryId"`
}

type UpdateTopicRequestHandler interface {
	Handle(ctx context.Context, req UpdateTopicRequest) (*user.Topic, error)
}

type updateTopicRequestHandler struct {
	repo user.Repository
}

func NewUpdateTopicHandler(repo user.Repository) UpdateTopicRequestHandler {
	return &updateTopicRequestHandler{
		repo: repo,
	}
}

func (h *updateTopicRequestHandler) Handle(ctx context.Context, req UpdateTopicRequest) (*user.Topic, error) {
	topic := &user.Topic{
		UserID:     req.User.ID,
		CategoryID: req.CategoryID,
		ID:         req.TopicID,
		Title:      req.Title,
		Content:    req.Content,
		ImagePath:  req.ImagePath,
	}

	err := h.repo.UpdateTopic(ctx, topic)
	if err != nil {
		return nil, err
	}
	return topic, nil
}
