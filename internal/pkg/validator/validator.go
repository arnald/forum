package validator

import (
	"fmt"
	"reflect"
	"regexp"
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
	Rules []func(interface{}) (bool, string)
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
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
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

func ValidateStruct(v *Validator, data interface{}, rules []ValidationRule) {
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

func required(value interface{}) (bool, string) {
	str, ok := value.(string)
	if !ok {
		return false, InvalidType
	}
	return str != "", "must be provided"
}

func minLength(minimumLenght int) func(interface{}) (bool, string) {
	return func(value interface{}) (bool, string) {
		str, ok := value.(string)
		if !ok {
			return false, InvalidType
		}
		return len(str) >= minimumLenght, fmt.Sprintf("must be at least %d characters long", minimumLenght)
	}
}

func maxLength(maximumLenght int) func(interface{}) (bool, string) {
	return func(value interface{}) (bool, string) {
		str, ok := value.(string)
		if !ok {
			return false, InvalidType
		}
		return len(str) <= maximumLenght, fmt.Sprintf("must be %d characters maximum", maximumLenght)
	}
}

func validEmail(value interface{}) (bool, string) {
	str, ok := value.(string)
	if !ok {
		return false, InvalidType
	}
	return Matches(str, EmailRX), InvalidEmail
}
