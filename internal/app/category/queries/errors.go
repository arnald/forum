package queries

import "errors"

var (
	ErrEmptyName            = errors.New("category name cannot be empty")
	ErrEmptyID              = errors.New("category ID cannot be empty")
	ErrCategoryNotFound     = errors.New("category not found")
	ErrCategoryAlreadyExists = errors.New("category with this name already exists")
	ErrCategoryHasPosts     = errors.New("cannot delete category that contains posts")
)