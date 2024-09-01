package validator

// Package for validating form inputs
import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

func ValidRating(value string) bool {
	switch strings.ToLower(value) {
	case
		"pg",
		"pg-13",
		"r":
		return true
	}
	return false
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}
