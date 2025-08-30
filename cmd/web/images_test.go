package main

import (
	"mime/multipart"
	"os"
	"path"
	"testing"

	"sketchdb.cozycole.net/internal/assert"
	"sketchdb.cozycole.net/internal/img"
	"sketchdb.cozycole.net/internal/utils"
)

type Size struct {
	Width  int
	Height int
}

func TestSaveLargeThumbnail(t *testing.T) {
	app := application{
		fileStorage: &img.FileStorage{RootPath: "/tmp"},
	}

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
				"medium": {Width: 320, Height: 180},
				"large":  {Width: 320, Height: 180},
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

func TestSaveMediumThumbnail(t *testing.T) {
	app := application{
		fileStorage: &img.FileStorage{RootPath: "/tmp"},
	}

	directoryName := "test-save-medium-thumb"

	thumbnail192x144, err := utils.CreateMultipartFileHeader("./testdata/test-thumbnail-192x144.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}

	thumbnail626x209, err := utils.CreateMultipartFileHeader("./testdata/test-thumbnail-626x209.jpg")
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
			name:          "192x144",
			thumbnail:     thumbnail192x144,
			thumbnailName: "192x144.jpg",
			desiredDimensions: map[string]Size{
				"small":  {Width: 320, Height: 180},
				"medium": {Width: 320, Height: 180},
			},
		},
		{
			name:          "626x209",
			thumbnail:     thumbnail626x209,
			thumbnailName: "626x209.jpg",
			desiredDimensions: map[string]Size{
				"small":  {Width: 320, Height: 180},
				"medium": {Width: 372, Height: 209},
			},
		},
		{
			name:          "320x240",
			thumbnail:     thumbnail320x240,
			thumbnailName: "320x240.jpg",
			desiredDimensions: map[string]Size{
				"small":  {Width: 320, Height: 180},
				"medium": {Width: 320, Height: 180},
			},
		},
		{
			name:          "800x450",
			thumbnail:     thumbnail800x450,
			thumbnailName: "800x450.jpg",
			desiredDimensions: map[string]Size{
				"small":  {Width: 320, Height: 180},
				"medium": {Width: 640, Height: 360},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app.saveMediumThumbnail(tt.thumbnailName, directoryName, tt.thumbnail)
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

func TestSaveLargeProfile(t *testing.T) {
	app := application{
		fileStorage: &img.FileStorage{RootPath: "/tmp"},
	}

	directoryName := "test-save-large-profile"

	thumbnail192x144, err := utils.CreateMultipartFileHeader("./testdata/test-thumbnail-192x144.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}

	thumbnail800x450, err := utils.CreateMultipartFileHeader("./testdata/test-thumbnail-800x450.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}

	thumbnail578x850, err := utils.CreateMultipartFileHeader("./testdata/test-profile-578x850.jpg")
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
			name:          "192x144",
			thumbnail:     thumbnail192x144,
			thumbnailName: "192x144.jpg",
			desiredDimensions: map[string]Size{
				"small":  {Width: 88, Height: 88},
				"medium": {Width: 144, Height: 144},
				"large":  {Width: 144, Height: 144},
			},
		},
		{
			name:          "800x450",
			thumbnail:     thumbnail800x450,
			thumbnailName: "800x450.jpg",
			desiredDimensions: map[string]Size{
				"small":  {Width: 88, Height: 88},
				"medium": {Width: 256, Height: 256},
				"large":  {Width: 450, Height: 450},
			},
		},
		{
			name:          "578x850.jpg",
			thumbnail:     thumbnail578x850,
			thumbnailName: "578x850.jpg.jpg",
			desiredDimensions: map[string]Size{
				"small":  {Width: 88, Height: 88},
				"medium": {Width: 256, Height: 256},
				"large":  {Width: 512, Height: 512},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app.saveLargeProfile(tt.thumbnailName, directoryName, tt.thumbnail)
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

func TestSaveMediumProfile(t *testing.T) {
	app := application{
		fileStorage: &img.FileStorage{RootPath: "/tmp"},
	}

	directoryName := "test-save-medium-profile"

	thumbnail192x144, err := utils.CreateMultipartFileHeader("./testdata/test-thumbnail-192x144.jpg")
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
			name:          "192x144",
			thumbnail:     thumbnail192x144,
			thumbnailName: "192x144.jpg",
			desiredDimensions: map[string]Size{
				"small":  {Width: 88, Height: 88},
				"medium": {Width: 144, Height: 144},
			},
		},
		{
			name:          "800x450",
			thumbnail:     thumbnail800x450,
			thumbnailName: "800x450.jpg",
			desiredDimensions: map[string]Size{
				"small":  {Width: 88, Height: 88},
				"medium": {Width: 256, Height: 256},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app.saveMediumProfile(tt.thumbnailName, directoryName, tt.thumbnail)
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
