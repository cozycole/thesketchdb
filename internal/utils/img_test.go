package utils

import (
	"image"
	"io"
	"os"
	"testing"
)

func TestCropImage(t *testing.T) {
	testImg := "./testdata/test-img.jpg"

	file, err := os.Open(testImg)
	if err != nil {
		t.Fatal(err)
	}

	rect := image.Rect(0, 45, 480, 315)
	buf, err := CropImg(file, ".jpg", rect)
	if err != nil {
		t.Fatal(err)
	}

	dst, err := os.Create("/tmp/test_output.jpg")
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.Copy(dst, buf)
	t.Run("Image created", func(t *testing.T) {
		if err != nil {
			t.Error(err)
		}
	})
}
