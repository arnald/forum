package oauthlogin

import "errors"

var (
	ErrInParameters  = errors.New("oauth: error in parameters")
	ErrStateMismatch = errors.New("oauth: state mismatch")
	ErrCodeMissing   = errors.New("oauth: code missing")
	ErrUserNotFound  = errors.New("oauth: user not found")
	ErrLoginFailed   = errors.New("oauth: login failed")
	ErrTokenInvalid  = errors.New("oauth: token invalid")
)
