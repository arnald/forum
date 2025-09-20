package queries

import "errors"

var (
	ErrEmptyTitle    = errors.New("title cannot be empty")
	ErrEmptyContent  = errors.New("content cannot be empty")
	ErrEmptyUserID   = errors.New("user ID cannot be empty")
	ErrPostNotFound  = errors.New("post not found")
	ErrEmptyCategory = errors.New("category cannot be empty")
)