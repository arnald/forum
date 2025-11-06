package validation

import (
	"net/mail"
	"strings"
	"unicode"
)

const (
	minUsernameLength = 3
	minPasswordLength = 8
)

// ValidateEmail validates the email format.
func ValidateEmail(email string) string {
	if email == "" {
		return "Email is required."
	}
	if !IsValidEmail(email) {
		return "Invalid email format."
	}
	return ""
}

// ValidateUsername validates the username length and presence.
func ValidateUsername(username string) string {
	if username == "" {
		return "Username is required."
	}
	if len(username) < minUsernameLength {
		return "Username must be at least 3 characters"
	}
	return ""
}

// IsValidEmail checks if the email format is valid.
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// HasLower checks if string contains lowercase letter.
func HasLower(s string) bool {
	for _, c := range s {
		if unicode.IsLower(c) {
			return true
		}
	}
	return false
}

// HasUpper checks if string contains uppercase letter.
func HasUpper(s string) bool {
	for _, c := range s {
		if unicode.IsUpper(c) {
			return true
		}
	}
	return false
}

// HasDigit checks if string contains digit.
func HasDigit(s string) bool {
	for _, c := range s {
		if unicode.IsDigit(c) {
			return true
		}
	}
	return false
}

// HasSpecial checks if string contains special character.
func HasSpecial(s string) bool {
	for _, c := range s {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && !unicode.IsSpace(c) {
			return true
		}
	}
	return false
}

// ValidatePassword validates password strength and requirements.
func ValidatePassword(pw string) string {
	if pw == "" {
		return "Password is required."
	}
	if len(pw) < minPasswordLength {
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
