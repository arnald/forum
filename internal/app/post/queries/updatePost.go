package queries

import (
	"context"
	"time"

	"github.com/arnald/forum/internal/domain/post"
)

type UpdatePostRequest struct {
	ID         string
	Title      string
	Content    string
	UserID     string
	Categories []string
}

type UpdatePostRequestHandler interface {
	Handle(ctx context.Context, req UpdatePostRequest) (*post.Post, error)
}

type updatePostRequestHandler struct {
	repo post.Repository
}

func NewUpdatePostHandler(repo post.Repository) UpdatePostRequestHandler {
	return &updatePostRequestHandler{
		repo: repo,
	}
}

func (h *updatePostRequestHandler) Handle(ctx context.Context, req UpdatePostRequest) (*post.Post, error) {
	if req.ID == "" {
		return nil, ErrPostNotFound
	}
	if req.Title == "" {
		return nil, ErrEmptyTitle
	}
	if req.Content == "" {
		return nil, ErrEmptyContent
	}

	existingPost, err := h.repo.GetPostByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if existingPost.UserID != req.UserID {
		return nil, ErrPostNotFound
	}

	updatedPost := &post.Post{
		ID:         req.ID,
		Title:      req.Title,
		Content:    req.Content,
		UserID:     req.UserID,
		Categories: req.Categories,
		CreatedAt:  existingPost.CreatedAt,
		UpdatedAt:  time.Now(),
	}

	err = h.repo.UpdatePost(ctx, updatedPost)
	if err != nil {
		return nil, err
	}

	return h.repo.GetPostByID(ctx, req.ID)
}