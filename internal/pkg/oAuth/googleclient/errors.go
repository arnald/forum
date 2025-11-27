package googleclient

import "errors"

var (
	ErrFailedToExchangeCode    = errors.New("oauth: failed to exchange code for token")
	ErrFailedToGetUser         = errors.New("oauth: failed to get user info")
	ErrFailedToParseToken      = errors.New("oauth: failed to parse token response")
	ErrFailedToParseUser       = errors.New("oauth: failed to parse user response")
	ErrTokenNotFound           = errors.New("oauth: token not found")
	ErrFailedToGetPrimaryEmail = errors.New("oauth: failed to get primary email")
	ErrNoVerifiedEmailFound    = errors.New("oauth: no verified email found")
	ErrEmailNotVerified        = errors.New("oauth: email not verified")
)
