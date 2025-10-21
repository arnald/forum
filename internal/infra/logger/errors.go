package logger

import "errors"

var (
	ErrInvalidRequestMethod  = errors.New("invalid request method")
	ErrInvalidRequestBody    = errors.New("invalid request body")
	ErrValidationFailed      = errors.New("validation failed")
	ErrUserNotFoundInContext = errors.New("user not found in context")
)
