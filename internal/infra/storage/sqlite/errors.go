package sqlite

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDuplicateUsername = errors.New("username already exists")
	ErrConstraint        = errors.New("sqlite constrain error")
	ErrUnknownConstraint = errors.New("sqlite unknow constraint error")
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrUserNotFound      = errors.New("user not found")
	ErrTopicNotFound     = errors.New("topic not found")
)

func MapSQLiteError(err error) error {
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		if sqliteErr.Code == sqlite3.ErrConstraint {
			msg := err.Error()

			switch {
			case strings.Contains(msg, "users.email"):
				return ErrDuplicateEmail
			case strings.Contains(msg, "users.username"):
				return ErrDuplicateUsername
			case strings.Contains(msg, "email LIKE"):
				return ErrInvalidEmail
			default:
				return fmt.Errorf("%w: %w", ErrConstraint, sqliteErr)
			}
		}
		return fmt.Errorf("%w: %w: %s ", ErrUnknownConstraint, sqliteErr.Code, sqliteErr.Error())
	}
	return err
}
