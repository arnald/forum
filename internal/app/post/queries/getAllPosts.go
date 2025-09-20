package queries

import (
	"context"

	"github.com/arnald/forum/internal/domain/post"
)

type GetAllPostsRequest struct{}

type GetAllPostsRequestHandler interface {
	Handle(ctx context.Context, req GetAllPostsRequest) ([]*post.Post, error)
}

type getAllPostsRequestHandler struct {
	repo post.Repository
}

func NewGetAllPostsHandler(repo post.Repository) GetAllPostsRequestHandler {
	return &getAllPostsRequestHandler{
		repo: repo,
	}
}

func (h *getAllPostsRequestHandler) Handle(ctx context.Context, req GetAllPostsRequest) ([]*post.Post, error) {
	return h.repo.GetAllPosts(ctx)
}