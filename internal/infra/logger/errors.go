package logger

import "errors"

var (
	ErrInvalidRequestMethod = errors.New("invalid request method")
	ErrInvalidRequestBody   = errors.New("invalid request body")
	ErrValidationFailed     = errors.New("validation failed")
)
