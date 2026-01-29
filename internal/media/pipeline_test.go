package media

import (
	"image"
	_ "image/jpeg"
	"os"
	"path"
	"testing"

	"sketchdb.cozycole.net/internal/assert"
	"sketchdb.cozycole.net/internal/fileStore"
)

func TestCreateImageVariantSpec(t *testing.T) {
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

	f2, err := os.Open("./testdata/test-thumbnail-626x209.jpg")
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

	tests := []struct {
		name          string
		img           image.Image
		thumbnailName string
		imgType       ImageType
		maxSize       Size
		variants      []VariantSpec
	}{
		{
			name:          "Thumbnail",
			img:           thumbnail,
			thumbnailName: "thumbnail.jpg",
			imgType:       Thumbnail,
			maxSize:       Large,
			variants: []VariantSpec{
				{
					Name:    "small",
					Width:   SmallThumbnailWidth,
					Height:  SmallThumbnailHeight,
					Mode:    FitCover,
					Quality: JPGQuality,
					Format:  FormatJPEG,
				},
				{
					Name:    "medium",
					Width:   MediumThumbnailWidth,
					Height:  MediumThumbnailHeight,
					Mode:    FitCover,
					Quality: JPGQuality,
					Format:  FormatJPEG,
				},
				{
					Name:    "large",
					Width:   LargeThumbnailWidth,
					Height:  LargeThumbnailHeight,
					Mode:    FitCover,
					Quality: JPGQuality,
					Format:  FormatJPEG,
				},
			},
		},
		{
			name:          "Profile",
			img:           profile,
			thumbnailName: "profile.jpg",
			imgType:       Profile,
			maxSize:       Medium,
			variants: []VariantSpec{
				{
					Name:    "small",
					Width:   SmallProfileWidth,
					Height:  SmallProfileWidth,
					Mode:    FitCover,
					Quality: JPGQuality,
					Format:  FormatJPEG,
				},
				{
					Name:    "medium",
					Width:   209,
					Height:  209,
					Mode:    FitCover,
					Quality: JPGQuality,
					Format:  FormatJPEG,
				},
			},
		},
		{
			name:          "Thumbnail Medium",
			img:           profile,
			thumbnailName: "thumbnail.jpg",
			imgType:       "thumbnail",
			maxSize:       Medium,
			variants: []VariantSpec{
				{
					Name:    "small",
					Width:   SmallThumbnailWidth,
					Height:  SmallThumbnailHeight,
					Mode:    FitCover,
					Quality: JPGQuality,
					Format:  FormatJPEG,
				},
				{
					Name:    "medium",
					Width:   372,
					Height:  209,
					Mode:    FitCover,
					Quality: JPGQuality,
					Format:  FormatJPEG,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			variants, err := createImageVariantSpec(tt.img, tt.maxSize, tt.imgType)
			assert.NilError(t, err)
			assert.Equal(t, len(variants), len(tt.variants))
			for i := range tt.variants {
				if variants[i] != tt.variants[i] {
					t.Errorf("mismatch at index %d: got %+v, want %+v", i, variants[i], tt.variants[i])
				}
			}
		})
	}
}

func TestSaveImageVariants(t *testing.T) {
	directory := path.Join(os.TempDir(), "test-save-image-variants")
	imgStorage := fileStore.FileStorage{RootPath: directory}

	err := os.MkdirAll(directory, 0o755)
	if err != nil {
		t.Fatal(err)
		return
	}

	inspectImage := false
	if !inspectImage {
		defer os.RemoveAll(directory)
	}

	f1, err := os.ReadFile("./testdata/test-thumbnail-626x209.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}

	// f2, err := os.ReadFile("./testdata/test-profile-578x850.jpg")
	// if err != nil {
	// 	t.Fatal(err)
	// 	return
	// }

	err = RunImagePipeline(
		f1,
		Medium,
		Thumbnail,
		"abcde.jpg",
		"/cast/thumbnail",
		&imgStorage,
	)

	if err != nil {
		t.Fatal(err)
		return
	}
}
