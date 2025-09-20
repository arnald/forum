package queries

import "errors"

var (
	ErrEmptyUserID   = errors.New("user ID cannot be empty")
	ErrEmptyTargetID = errors.New("target ID cannot be empty")
	ErrVoteNotFound  = errors.New("vote not found")
)