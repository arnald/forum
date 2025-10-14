package topics

import "errors"

var (
	ErrTopicNotFound = errors.New("topic not found")
	ErrUserNotFound  = errors.New("user not found")
)
