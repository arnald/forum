package validation

import (
	"net/mail"
	"strings"
	"unicode"
)

func ValidateEmail(email string) string {
	if email == "" {
		return "Email is required."
	}
	if !IsValidEmail(email) {
		return "Invalid email format."
	}
	return ""
}

func ValidateUsername(username string) string {
	if username == "" {
		return "Username is required."
	}
	if len(username) < 3 {
		return "Username must be at least 3 characters"
	}
	return ""
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func HasLower(s string) bool {
	for _, c := range s {
		if unicode.IsLower(c) {
			return true
		}
	}
	return false
}

func HasUpper(s string) bool {
	for _, c := range s {
		if unicode.IsUpper(c) {
			return true
		}
	}
	return false
}

func HasDigit(s string) bool {
	for _, c := range s {
		if unicode.IsDigit(c) {
			return true
		}
	}
	return false
}

func HasSpecial(s string) bool {
	for _, c := range s {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && !unicode.IsSpace(c) {
			return true
		}
	}
	return false
}

func ValidatePassword(pw string) string {
	if pw == "" {
		return "Password is required."
	}
	if len(pw) < 8 {
		return "Password must be 8+ chars"
	}

	var missing []string
	if !HasLower(pw) {
		missing = append(missing, "lowercase letter")
	}
	if !HasUpper(pw) {
		missing = append(missing, "uppercase letter")
	}
	if !HasDigit(pw) {
		missing = append(missing, "number")
	}
	if !HasSpecial(pw) {
		missing = append(missing, "special character")
	}

	if len(missing) > 0 {
		return "Missing: " + strings.Join(missing, ", ")
	}

	return "" // empty string means no error
}
