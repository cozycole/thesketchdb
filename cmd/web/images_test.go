package main

import (
	"mime/multipart"
	"os"
	"path"
	"testing"

	"sketchdb.cozycole.net/internal/assert"
	"sketchdb.cozycole.net/internal/utils"
)

type Size struct {
	Width  int
	Height int
}

func TestSaveLargeThumbnail(t *testing.T) {
	app := newTestApplication(t)

	directoryName := "test-save-large-thumb"

	// Each of these files test a different logical output for saving files based on their size
	thumbnail1920x1080, err := utils.CreateMultipartFileHeader("./testdata/test-thumbnail-1920x1080.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}

	thumbnail320x240, err := utils.CreateMultipartFileHeader("./testdata/test-thumbnail-320x240.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}

	thumbnail800x450, err := utils.CreateMultipartFileHeader("./testdata/test-thumbnail-800x450.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}

	tests := []struct {
		name              string
		thumbnail         *multipart.FileHeader
		thumbnailName     string
		desiredDimensions map[string]Size
	}{
		{
			name:          "1920x1080",
			thumbnail:     thumbnail1920x1080,
			thumbnailName: "1920x1080.jpg",
			desiredDimensions: map[string]Size{
				"small":  {Width: 320, Height: 180},
				"medium": {Width: 640, Height: 360},
				"large":  {Width: 1280, Height: 720},
			},
		},
		{
			name:          "320x240",
			thumbnail:     thumbnail320x240,
			thumbnailName: "320x240.jpg",
			desiredDimensions: map[string]Size{
				"small":  {Width: 320, Height: 180},
				"medium": {Width: 640, Height: 360},
				"large":  {Width: 640, Height: 360},
			},
		},
		{
			name:          "800x450",
			thumbnail:     thumbnail800x450,
			thumbnailName: "800x450.jpg",
			desiredDimensions: map[string]Size{
				"small":  {Width: 320, Height: 180},
				"medium": {Width: 640, Height: 360},
				"large":  {Width: 800, Height: 450},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app.saveLargeThumbnail(tt.thumbnailName, directoryName, tt.thumbnail)
			if err != nil {
				t.Fatal(err)
			}
			for size, dimensions := range tt.desiredDimensions {
				thumbnail, err := os.Open(path.Join(os.TempDir(), directoryName, size, tt.thumbnailName))
				if err != nil {
					t.Fatal(err)
				}

				width, height, err := utils.GetImageDimensions(thumbnail)
				assert.Equal(t, width, dimensions.Width)
				assert.Equal(t, height, dimensions.Height)
			}
		})
	}
}
