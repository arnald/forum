package oauthrepo

import "errors"

var (
	ErrUserNotFound              = errors.New("oauth user not found")
	ErrTransactionRollbackFailed = errors.New("transaction rollback failed")
	ErrTransactionCommitFailed   = errors.New("transaction commit failed")
)
