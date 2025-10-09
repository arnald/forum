package categories

import "errors"

var (
	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrCategoryNotFound      = errors.New("category not found")
	ErrUserNotFound          = errors.New("user not found")
)
