package commentcommands

import (
	"context"

	"github.com/arnald/forum/internal/domain/comment"
	"github.com/arnald/forum/internal/domain/user"
)

type DeleteCommentRequest struct {
	User      *user.User
	CommentID int `json:"commentId"`
}

type DeleteCommentRequestHandler interface {
	Handle(ctx context.Context, req DeleteCommentRequest) error
}

type deleteCommentRequestHandler struct {
	repo comment.Repository
}

func NewDeleteCommentHandler(repo comment.Repository) DeleteCommentRequestHandler {
	return &deleteCommentRequestHandler{
		repo: repo,
	}
}

func (h *deleteCommentRequestHandler) Handle(ctx context.Context, req DeleteCommentRequest) error {
	err := h.repo.DeleteComment(ctx, req.User.ID, req.CommentID)
	if err != nil {
		return err
	}

	return nil
}
