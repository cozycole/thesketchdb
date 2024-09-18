package validator

// Package for validating form inputs
import (
	"mime/multipart"
	"net/http"
	"slices"
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
	_, err := time.Parse(time.DateOnly, value)
	return err == nil
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func IsMime(file multipart.File, mtypes ...string) bool {
	buf := make([]byte, 512)

	if _, err := file.Read(buf); err != nil {
		return false
	}
	return slices.Contains(mtypes, http.DetectContentType(buf))
}

func BoolWithError(input bool, err error) bool {
	if err != nil {
		return false
	}
	return input
}

func IsZero(num int) bool {
	return num == 0
}
