package comments

import "errors"

var (
	ErrCommentNotFound = errors.New("comment not found")
	ErrFailedToUpdate  = errors.New("failed to update comment, user not found or not authorized")
)

// ErrTopicNotFound   = errors.New("topic not found")
// ErrUserNotFound    = errors.New("user not found").
