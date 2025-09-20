package queries

import (
	"context"
	"time"

	"github.com/arnald/forum/internal/domain/comment"
	"github.com/arnald/forum/internal/pkg/uuid"
)

type CreateCommentRequest struct {
	Content  string
	PostID   string
	UserID   string
	ParentID *string // Optional for replies
}

type CreateCommentRequestHandler interface {
	Handle(ctx context.Context, req CreateCommentRequest) (*comment.Comment, error)
}

type createCommentRequestHandler struct {
	repo         comment.Repository
	uuidProvider uuid.Provider
}

func NewCreateCommentHandler(repo comment.Repository, uuidProvider uuid.Provider) CreateCommentRequestHandler {
	return &createCommentRequestHandler{
		repo:         repo,
		uuidProvider: uuidProvider,
	}
}

func (h *createCommentRequestHandler) Handle(ctx context.Context, req CreateCommentRequest) (*comment.Comment, error) {
	if req.Content == "" {
		return nil, ErrEmptyContent
	}
	if req.PostID == "" {
		return nil, ErrEmptyPostID
	}
	if req.UserID == "" {
		return nil, ErrEmptyUserID
	}

	level := 0
	if req.ParentID != nil {
		parentComment, err := h.repo.GetCommentByID(ctx, *req.ParentID)
		if err != nil {
			return nil, ErrParentCommentNotFound
		}
		level = parentComment.Level + 1

		// Limit nesting depth
		if level > maxNestingLevel {
			return nil, ErrMaxNestingExceeded
		}
	}

	newComment := &comment.Comment{
		ID:        h.uuidProvider.NewUUID(),
		Content:   req.Content,
		PostID:    req.PostID,
		UserID:    req.UserID,
		ParentID:  req.ParentID,
		Level:     level,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := h.repo.CreateComment(ctx, newComment)
	if err != nil {
		return nil, err
	}

	return newComment, nil
}