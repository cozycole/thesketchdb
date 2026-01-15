package wikipedia_test

import (
	"testing"

	"sketchdb.cozycole.net/internal/services/wikipedia"
)

func TestGetExtract(t *testing.T) {
	extract, err := wikipedia.GetExtract("Shane_Gillis")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if extract == "" {
		t.Fatal("expected a non-empty extract, got empty string")
	}

	// Log a preview for human verification
	if len(extract) > 100 {
		t.Logf("Extract (first 100 chars): %s...", extract[:100])
	} else {
		t.Logf("Extract: %s", extract)
	}
}
