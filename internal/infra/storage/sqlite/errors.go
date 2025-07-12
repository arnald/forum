package sqlite

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mattn/go-sqlite3"
)

var ErrDuplicateEmail = errors.New("duplicate email")
var ErrDuplicateUsername = errors.New("duplicate username")

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
			default:
				return fmt.Errorf("sqlite constraint error: %v", sqliteErr)
			}
		}
		return fmt.Errorf("sqlite error %d: %s", sqliteErr.Code, sqliteErr.Error())
	}
	return err
}
