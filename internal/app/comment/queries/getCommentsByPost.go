package queries

import (
	"context"

	"github.com/arnald/forum/internal/domain/comment"
)

type GetCommentsByPostRequest struct {
	PostID string
}

type GetCommentsByPostRequestHandler interface {
	Handle(ctx context.Context, req GetCommentsByPostRequest) ([]*comment.Comment, error)
}

type getCommentsByPostRequestHandler struct {
	repo comment.Repository
}

func NewGetCommentsByPostHandler(repo comment.Repository) GetCommentsByPostRequestHandler {
	return &getCommentsByPostRequestHandler{
		repo: repo,
	}
}

func (h *getCommentsByPostRequestHandler) Handle(ctx context.Context, req GetCommentsByPostRequest) ([]*comment.Comment, error) {
	if req.PostID == "" {
		return nil, ErrEmptyPostID
	}

	return h.repo.GetCommentsByPostID(ctx, req.PostID)
}