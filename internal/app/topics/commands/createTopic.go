package topiccommands

import (
	"context"

	"github.com/arnald/forum/internal/domain/user"
)

type CreateTopicRequest struct {
	User       *user.User
	Title      string `json:"title"`
	Content    string `json:"content"`
	ImagePath  string `json:"imagePath"`
	CategoryID int    `json:"categoryId"`
}

type CreateTopicRequestHandler interface {
	Handle(ctx context.Context, req CreateTopicRequest) (*user.Topic, error)
}

type createTopicRequestHandler struct {
	repo user.Repository
}

func NewCreateTopicHandler(repo user.Repository) CreateTopicRequestHandler {
	return &createTopicRequestHandler{
		repo: repo,
	}
}

func (h *createTopicRequestHandler) Handle(ctx context.Context, req CreateTopicRequest) (*user.Topic, error) {
	topic := &user.Topic{
		UserID:     req.User.ID,
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Content:    req.Content,
		ImagePath:  req.ImagePath,
	}

	err := h.repo.CreateTopic(ctx, topic)
	if err != nil {
		return nil, err
	}
	return topic, nil
}
