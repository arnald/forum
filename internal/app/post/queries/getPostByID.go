package queries

import (
	"context"

	"github.com/arnald/forum/internal/domain/post"
)

type GetPostByIDRequest struct {
	ID string
}

type GetPostByIDRequestHandler interface {
	Handle(ctx context.Context, req GetPostByIDRequest) (*post.Post, error)
}

type getPostByIDRequestHandler struct {
	repo post.Repository
}

func NewGetPostByIDHandler(repo post.Repository) GetPostByIDRequestHandler {
	return &getPostByIDRequestHandler{
		repo: repo,
	}
}

func (h *getPostByIDRequestHandler) Handle(ctx context.Context, req GetPostByIDRequest) (*post.Post, error) {
	if req.ID == "" {
		return nil, ErrPostNotFound
	}

	return h.repo.GetPostByID(ctx, req.ID)
}