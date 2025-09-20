package queries

import "errors"

const maxNestingLevel = 5 // Maximum comment nesting depth

var (
	ErrEmptyContent           = errors.New("content cannot be empty")
	ErrEmptyPostID            = errors.New("post ID cannot be empty")
	ErrEmptyUserID            = errors.New("user ID cannot be empty")
	ErrCommentNotFound        = errors.New("comment not found")
	ErrParentCommentNotFound  = errors.New("parent comment not found")
	ErrMaxNestingExceeded     = errors.New("maximum comment nesting level exceeded")
)