package comments

import "errors"

var (
	ErrCommentNotFound = errors.New("comment not found")
	ErrTopicNotFound   = errors.New("topic not found")
	ErrUserNotFound    = errors.New("user not found")
)
