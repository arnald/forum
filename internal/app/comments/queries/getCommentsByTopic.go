package commentqueries

import (
	"context"

	"github.com/arnald/forum/internal/domain/comment"
)

type GetCommentsByTopicRequest struct {
	TopicID int `json:"topicId"`
}

type GetCommentsByTopicRequestHandler interface {
	Handle(ctx context.Context, req GetCommentsByTopicRequest) ([]comment.Comment, error)
}

type getCommentsByTopicRequestHandler struct {
	repo comment.Repository
}

func NewGetCommentsByTopicRequestHandler(repo comment.Repository) GetCommentsByTopicRequestHandler {
	return &getCommentsByTopicRequestHandler{
		repo: repo,
	}
}

func (h *getCommentsByTopicRequestHandler) Handle(ctx context.Context, req GetCommentsByTopicRequest) ([]comment.Comment, error) {
	comments, err := h.repo.GetCommentsByTopicID(ctx, req.TopicID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}
