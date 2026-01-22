package media

import (
	"image"
	_ "image/jpeg"
	"os"
	"path"
	"testing"

	"sketchdb.cozycole.net/internal/assert"
)

func TestProcessCoverImage(t *testing.T) {
	directory := path.Join(os.TempDir(), "test-process-image")
	err := os.MkdirAll(directory, 0o755)
	if err != nil {
		t.Fatal(err)
		return
	}

	inspectImage := false
	if !inspectImage {
		defer os.RemoveAll(directory)
	}

	f1, err := os.Open("./testdata/test-thumbnail-1920x1080.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}
	defer f1.Close()

	thumbnail, _, err := image.Decode(f1)
	if err != nil {
		t.Fatal(err)
		return
	}
	// thumbnail = RemoveBorders(thumbnail, 15, 2)

	f2, err := os.Open("./testdata/test-profile-578x850.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}
	defer f2.Close()

	profile, _, err := image.Decode(f2)
	if err != nil {
		t.Fatal(err)
		return
	}
	// profile = RemoveBorders(profile, 15, 2)

	tests := []struct {
		name          string
		img           image.Image
		thumbnailName string
		spec          VariantSpec
	}{
		{
			name:          "Thumbnail",
			img:           thumbnail,
			thumbnailName: "thumbnail.jpg",
			spec: VariantSpec{
				Width:   640,
				Height:  360,
				Mode:    FitCover,
				Format:  FormatJPEG,
				Quality: 85,
			},
		},
		{
			name:          "Profile",
			img:           profile,
			thumbnailName: "profile.jpg",
			spec: VariantSpec{
				Width:   256,
				Height:  256,
				Mode:    FitCover,
				Format:  FormatJPEG,
				Quality: 85,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := Process(tt.img, []VariantSpec{tt.spec})
			assert.NilError(t, err)
			if len(out) == 0 {
				t.Errorf("%s image has empty output", tt.name)
			}

			err = os.WriteFile(
				path.Join(directory, tt.thumbnailName),
				out[0].Bytes,
				0o644,
			)

			assert.NilError(t, err)
		})
	}
}
