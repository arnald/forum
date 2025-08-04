package helpers

// import (
// 	"errors"
// 	"strings"
// )

// var (
// 	ErrEmptyUsername = errors.New("empty username").Error()
// 	ErrEmptyPassword = errors.New("empty password").Error()
// )

// func UserRegisterIsValid(usename, password, email string) error {
// 	var errorBuilder strings.Builder
// 	if strings.TrimSpace(usename) == "" {
// 		errorBuilder.Write([]byte(ErrEmptyUsername))
// 		errorBuilder.Write([]byte(","))
// 	}
// 	if strings.TrimSpace(password) == "" {
// 		errorBuilder.Write([]byte(" "))
// 		errorBuilder.Write([]byte(ErrEmptyPassword))
// 		errorBuilder.Write([]byte(","))
// 	}

// 	err := ValidateEmail(email)
// 	if err != nil {
// 		errorBuilder.Write([]byte(" "))
// 		errorBuilder.Write([]byte(err.Error()))
// 	}

// 	errorStr := strings.TrimSpace(errorBuilder.String())
// 	if errorStr == "" {
// 		return nil
// 	}

// 	if errorStr[:len(errorStr)-1] == "," {
// 		errorStr = errorStr[:len(errorStr)-1]
// 	}

// 	return errors.New(errorStr)
// }
