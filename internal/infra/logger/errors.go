package logger

import "errors"

var (
	ErrInvalidRequestMethod  = errors.New("invalid request method")
	ErrInvalidRequestBody    = errors.New("invalid request body")
	ErrValidationFailed      = errors.New("validation failed")
	ErrUserNotFoundInContext = errors.New("user not found in context")
	ErrNeitherIDProvided     = errors.New("neither topic_id nor comment_id provided")
	ErrBothIDsProvided       = errors.New("both topic_id and comment_id provided")
)
