package validator

// Package for validating form inputs
import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = map[string]string{}
	}
	fmt.Print(v.FieldErrors[key])
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	fmt.Print(strings.TrimSpace(value) != "")
	return strings.TrimSpace(value) != ""
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

func ValidDate(value string) bool {
	_, err := time.Parse("2006-02-01", value)
	return err == nil
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}
