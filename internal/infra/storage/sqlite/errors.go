package sqlite

import "errors"

const uniqueConstraintViolationErrorCode = 1062

var ErrDuplicateEmail = errors.New("duplicate email")
