package comment

import "context"

type Repository interface {
	CreateComment(ctx context.Context, comment *Comment) error
	UpdateComment(ctx context.Context, comment *Comment) error
	DeleteComment(ctx context.Context, userID string, commentID int) error
	GetCommentByID(ctx context.Context, commentID int) (*Comment, error)
	GetCommentsByTopicID(ctx context.Context, topicID int) ([]Comment, error)
}
