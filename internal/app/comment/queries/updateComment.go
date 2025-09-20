package queries

import (
	"context"
	"time"

	"github.com/arnald/forum/internal/domain/comment"
)

type UpdateCommentRequest struct {
	ID      string
	Content string
	UserID  string
}

type UpdateCommentRequestHandler interface {
	Handle(ctx context.Context, req UpdateCommentRequest) (*comment.Comment, error)
}

type updateCommentRequestHandler struct {
	repo comment.Repository
}

func NewUpdateCommentHandler(repo comment.Repository) UpdateCommentRequestHandler {
	return &updateCommentRequestHandler{
		repo: repo,
	}
}

func (h *updateCommentRequestHandler) Handle(ctx context.Context, req UpdateCommentRequest) (*comment.Comment, error) {
	if req.ID == "" {
		return nil, ErrCommentNotFound
	}
	if req.Content == "" {
		return nil, ErrEmptyContent
	}
	if req.UserID == "" {
		return nil, ErrEmptyUserID
	}

	existingComment, err := h.repo.GetCommentByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if existingComment.UserID != req.UserID {
		return nil, ErrCommentNotFound
	}

	updatedComment := &comment.Comment{
		ID:        req.ID,
		Content:   req.Content,
		PostID:    existingComment.PostID,
		UserID:    req.UserID,
		ParentID:  existingComment.ParentID,
		Level:     existingComment.Level,
		CreatedAt: existingComment.CreatedAt,
		UpdatedAt: time.Now(),
	}

	err = h.repo.UpdateComment(ctx, updatedComment)
	if err != nil {
		return nil, err
	}

	return h.repo.GetCommentByID(ctx, req.ID)
}