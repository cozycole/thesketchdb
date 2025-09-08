package utils

import ()

// Helper functions for tests across several packages

func ptr[T any](v T) *T {
	return &v
}
