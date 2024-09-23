package models

import (
	"testing"
)

func TestCreateSlugName(t *testing.T) {
	maxLength := 40
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "SpecialCharacters",
			input: "Cole's Comedy Sketch #awesome",
			want:  "coles-comedy-sketch-awesome",
		},
		{
			name:  "EmptyString",
			input: "",
			want:  "",
		},
		{
			name:  "OverMaxLength",
			input: "This is a test This is a test This a test",
			want:  "this-is-a-test-this-is-a-test-this-a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := CreateSlugName(tt.input, maxLength)

			if output != tt.want {
				t.Errorf("got %q; want %q", output, tt.want)
			}
		})
	}
}
