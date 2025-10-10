package commentqueries

import (
	"context"

	"github.com/arnald/forum/internal/domain/comment"
)

type GetCommentRequest struct {
	CommentID int `json:"commentId"`
}

type GetCommentRequestHandler interface {
	Handle(ctx context.Context, req GetCommentRequest) (*comment.Comment, error)
}

type getCommentRequestHandler struct {
	repo comment.Repository
}

func NewGetCommentHandler(repo comment.Repository) GetCommentRequestHandler {
	return &getCommentRequestHandler{
		repo: repo,
	}
}

func (h *getCommentRequestHandler) Handle(ctx context.Context, req GetCommentRequest) (*comment.Comment, error) {
	comment, err := h.repo.GetCommentByID(ctx, req.CommentID)
	if err != nil {
		return nil, err
	}

	return comment, nil
}
