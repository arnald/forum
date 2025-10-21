package commentcommands

import (
	"context"

	"github.com/arnald/forum/internal/domain/comment"
	"github.com/arnald/forum/internal/domain/user"
)

type CreateCommentRequest struct {
	User    *user.User
	Content string `json:"content"`
	TopicID int    `json:"topicId"`
}

type CreateCommentRequestHandler interface {
	Handle(ctx context.Context, req CreateCommentRequest) (*comment.Comment, error)
}

type createCommentRequestHandler struct {
	repo comment.Repository
}

func NewCreateCommentRequestHandler(repo comment.Repository) CreateCommentRequestHandler {
	return &createCommentRequestHandler{
		repo: repo,
	}
}

func (h *createCommentRequestHandler) Handle(ctx context.Context, req CreateCommentRequest) (*comment.Comment, error) {
	comment := &comment.Comment{
		UserID:  req.User.ID,
		TopicID: req.TopicID,
		Content: req.Content,
	}

	err := h.repo.CreateComment(ctx, comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}
