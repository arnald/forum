package queries

import (
	"context"
	"time"

	"github.com/arnald/forum/internal/domain/post"
	"github.com/arnald/forum/internal/pkg/uuid"
)

type CreatePostRequest struct {
	Title      string
	Content    string
	UserID     string
	Categories []string
}

type CreatePostRequestHandler interface {
	Handle(ctx context.Context, req CreatePostRequest) (*post.Post, error)
}

type createPostRequestHandler struct {
	repo         post.Repository
	uuidProvider uuid.Provider
}

func NewCreatePostHandler(repo post.Repository, uuidProvider uuid.Provider) CreatePostRequestHandler {
	return &createPostRequestHandler{
		repo:         repo,
		uuidProvider: uuidProvider,
	}
}

func (h *createPostRequestHandler) Handle(ctx context.Context, req CreatePostRequest) (*post.Post, error) {
	if req.Title == "" {
		return nil, ErrEmptyTitle
	}
	if req.Content == "" {
		return nil, ErrEmptyContent
	}
	if req.UserID == "" {
		return nil, ErrEmptyUserID
	}

	newPost := &post.Post{
		ID:         h.uuidProvider.NewUUID(),
		Title:      req.Title,
		Content:    req.Content,
		UserID:     req.UserID,
		Categories: req.Categories,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := h.repo.CreatePost(ctx, newPost)
	if err != nil {
		return nil, err
	}

	return newPost, nil
}