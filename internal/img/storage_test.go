package img

import (
	"os"
	"path"
	"testing"

	"sketchdb.cozycole.net/internal/utils"
)

func TestSaveMultipartFile(t *testing.T) {
	storage := FileStorage{Path: os.TempDir()}

	header, err := utils.CreateMultipartFileHeader("./testdata/test-img.jpg")
	if err != nil {
		t.Fatal(err)
	}

	multipartFile, err := header.Open()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Valid storage", func(t *testing.T) {
		err = storage.SaveMultipartFile("test-img.jpg", multipartFile)
		if err != nil {
			t.Error(err)
		}
		_, err := os.Open(path.Join(storage.Path, "test-img.jpg"))
		if err != nil {
			t.Error(err)
		}
		os.Remove(path.Join(storage.Path, "test-img.jpg"))
	})

}
