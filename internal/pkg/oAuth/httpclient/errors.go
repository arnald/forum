package httpclient

import "errors"

var (
	ErrFailedToExecuteRequest   = errors.New("failed to execute request")
	ErrFailedToGetUser          = errors.New("failed to get user")
	ErrFailedToParseUser        = errors.New("failed to parse user")
	ErrFailedToGetPrimaryEmail  = errors.New("failed to get primary email")
	ErrFailedToReadResponseBody = errors.New("failed to read response body")
	ErrRequestFailedWithStatus  = errors.New("request failed with status")
)
