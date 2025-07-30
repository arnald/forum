package helpers

import (
	"errors"
	"regexp"
)

var (
	ErrEmptyEmail         = errors.New("empty email")
	ErrInvalidEmailFormat = errors.New("invalid email format")
)

const emailRegex = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

func ValidateEmail(s string) error {
	if s == "" {
		return ErrEmptyEmail
	}

	emailRegex := regexp.MustCompile(emailRegex)

	if !emailRegex.MatchString(s) {
		return ErrInvalidEmailFormat
	}

	return nil
}
