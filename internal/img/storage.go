package img

import (
	"errors"
	"io"
	"os"
	"path"
)

type FileStorageInterface interface {
	SaveFile(string, io.Reader) error
}

type FileStorage struct {
	Path string
}

func (s *FileStorage) SaveFile(subPath string, file io.Reader) error {
	imgPath := path.Join(s.Path, subPath)
	imgDir := path.Dir(imgPath)
	// Make the dir if it doesn't exist
	if _, err := os.Stat(imgPath); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(imgDir, 0755)
		if err != nil {
			return err
		}
	}

	dst, err := os.Create(imgPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return err
	}
	return nil
}
