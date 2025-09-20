package comment

import (
	"context"
)

type Repository interface {
	CreateComment(ctx context.Context, comment *Comment) error
	GetCommentByID(ctx context.Context, id string) (*Comment, error)
	GetCommentsByPostID(ctx context.Context, postID string) ([]*Comment, error)
	GetCommentsByUserID(ctx context.Context, userID string) ([]*Comment, error)
	GetCommentReplies(ctx context.Context, parentID string) ([]*Comment, error)
	UpdateComment(ctx context.Context, comment *Comment) error
	DeleteComment(ctx context.Context, id string) error
	GetCommentTree(ctx context.Context, postID string) ([]*CommentTree, error)
}