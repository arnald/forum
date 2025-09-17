package validator

import (
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

const (
	InvalidType  = "invalid type"
	InvalidEmail = "invalid email"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	Errors map[string]string
	Rules  map[string][]ValidationRule
}

type ValidationRule struct {
	Field string
	Rules []func(any) (bool, string)
}

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
		Rules:  make(map[string][]ValidationRule),
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// For future checks.
func In(value string, list ...string) bool {
	return slices.Contains(list, value)
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}

func ValidateStruct(v *Validator, data any, rules []ValidationRule) {
	val := reflect.ValueOf(data).Elem()

	for _, rule := range rules {
		field := val.FieldByName(rule.Field)
		if !field.IsValid() {
			continue
		}

		for _, validationFunc := range rule.Rules {
			ok, message := validationFunc(field.Interface())
			v.Check(ok, rule.Field, message)
		}
	}
}

func required(value any) (bool, string) {
	switch v := value.(type) {
	case string:
		return v != "", "must be provided"
	case int:
		return true, ""
	default:
		return false, InvalidType
	}
}

func optional(validationFunc func(any) (bool, string)) func(any) (bool, string) {
	return func(value any) (bool, string) {
		str, ok := value.(string)
		if !ok {
			return false, "field must be a string"
		}

		if str == "" {
			return true, ""
		}

		return validationFunc(value)
	}
}

func minLength(minimumLenght int) func(any) (bool, string) {
	return func(value any) (bool, string) {
		str, ok := value.(string)
		if !ok {
			return false, InvalidType
		}
		return len(str) >= minimumLenght, fmt.Sprintf("must be at least %d characters long", minimumLenght)
	}
}

func maxLength(maximumLenght int) func(any) (bool, string) {
	return func(value any) (bool, string) {
		str, ok := value.(string)
		if !ok {
			return false, InvalidType
		}
		return len(str) <= maximumLenght, fmt.Sprintf("must be %d characters maximum", maximumLenght)
	}
}

func isPositiveInt(value any) (bool, string) {
	num, ok := value.(int)
	if !ok {
		return false, InvalidType
	}
	return num > 0, "must be a positive integer"
}

func maxInt(limit int) func(any) (bool, string) {
	return func(value any) (bool, string) {
		num, ok := value.(int)
		if !ok {
			return false, InvalidType
		}
		return num <= limit, fmt.Sprintf("must be less than or equal to %d", limit)
	}
}

func validEmail(value any) (bool, string) {
	str, ok := value.(string)
	if !ok {
		return false, InvalidType
	}
	return Matches(str, EmailRX), InvalidEmail
}

func (v *Validator) ToStringErrors() string {
	strError := ""
	for key, value := range v.Errors {
		strError += key + ": " + value + " "
	}
	return strings.TrimSpace(strError)
}

func validImagePath(value any) (bool, string) {
	validImageExtensions := map[string]bool{
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".gif":  true,
	}

	str, ok := value.(string)
	if !ok {
		return false, InvalidType
	}
	ext := strings.ToLower(filepath.Ext(str))
	return validImageExtensions[ext], "must be a valid image file"
}

// var validCategories = map[string]bool{
// 	"General Discussion": true,
// 	"Feedback":           true,
// 	"Off-Topic":          true,
// }

// func validCategory(value any) (bool, string) {
// 	str, ok := value.(string)
// 	if !ok {
// 		return false, InvalidType
// 	}
// 	return validCategories[str], "must be a valid category"
// }

func validOrderBy(value any) (bool, string) {
	orderByWhitelist := map[string]bool{
		"created_at ASC":  true,
		"created_at DESC": true,
		"updated_at ASC":  true,
		"updated_at DESC": true,
		"title ASC":       true,
		"title DESC":      true,
	}

	str, ok := value.(string)
	if !ok {
		return false, InvalidType
	}
	if str == "" {
		return true, ""
	}
	return orderByWhitelist[str], "must be a valid order by field"
}
