package commentcommands

import (
	"context"

	"github.com/arnald/forum/internal/domain/comment"
	"github.com/arnald/forum/internal/domain/user"
)

type UpdateCommentRequest struct {
	User      *user.User
	Content   string `json:"content"`
	CommentID int    `json:"commentId"`
}

type UpdateCommentRequestHandler interface {
	Handle(ctx context.Context, req UpdateCommentRequest) (*comment.Comment, error)
}

type updateCommentRequestHandler struct {
	repo comment.Repository
}

func NewUpdateCommentRequestHandler(repo comment.Repository) UpdateCommentRequestHandler {
	return &updateCommentRequestHandler{
		repo: repo,
	}
}

func (h *updateCommentRequestHandler) Handle(ctx context.Context, req UpdateCommentRequest) (*comment.Comment, error) {
	comment := &comment.Comment{
		ID:      req.CommentID,
		UserID:  req.User.ID,
		Content: req.Content,
	}

	err := h.repo.UpdateComment(ctx, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}
