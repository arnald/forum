package comment

import "context"

type Repository interface {
	CreateComment(ctx context.Context, comment *Comment) error
	UpdateComment(ctx context.Context, comment *Comment) error
	DeleteComment(ctx context.Context, userID string, commentID int) error
	GetCommentByID(ctx context.Context, commentID int) (*Comment, error)      // TODO: make it return votes
	GetCommentsByTopicID(ctx context.Context, topicID int) ([]Comment, error) // TODO: clean up (not returning votes)
	GetCommentsWithVotes(ctx context.Context, topicID int, userID *string) ([]Comment, error)
}
