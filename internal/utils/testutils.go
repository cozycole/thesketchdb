package utils

import (
	"time"
)

// Helper functions for tests across several packages

func GetStringPtr(value string) *string {
	s := value
	return &s
}

func GetIntPtr(value int) *int {
	i := value
	return &i
}

func GetTimePtr(time time.Time) *time.Time {
	t := time
	return &t
}
