package queries

import (
	"context"

	"github.com/arnald/forum/internal/domain/comment"
)

type GetCommentTreeRequest struct {
	PostID string
}

type GetCommentTreeRequestHandler interface {
	Handle(ctx context.Context, req GetCommentTreeRequest) ([]*comment.CommentTree, error)
}

type getCommentTreeRequestHandler struct {
	repo comment.Repository
}

func NewGetCommentTreeHandler(repo comment.Repository) GetCommentTreeRequestHandler {
	return &getCommentTreeRequestHandler{
		repo: repo,
	}
}

func (h *getCommentTreeRequestHandler) Handle(ctx context.Context, req GetCommentTreeRequest) ([]*comment.CommentTree, error) {
	if req.PostID == "" {
		return nil, ErrEmptyPostID
	}

	return h.repo.GetCommentTree(ctx, req.PostID)
}