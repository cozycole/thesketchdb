package utils

import (
	"testing"
)

// Helper functions for tests across several packages

func ptr[T any](v T) *T {
	return &v
}

func AssertCallCount(
	t *testing.T,
	method string,
	got int,
	shouldCall bool,
) {
	t.Helper()

	want := 0
	if shouldCall {
		want = 1
	}

	if got != want {
		t.Errorf(
			"%s called %d times, want %d",
			method,
			got,
			want,
		)
	}
}
