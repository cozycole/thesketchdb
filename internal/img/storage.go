package img

import (
	"errors"
	"io"
	"os"
	"path"
)

type FileStorageInterface interface {
	SaveFile(string, io.Reader) error
	DeleteFile(string) error
}

type FileStorage struct {
	RootPath string
}

func (s *FileStorage) SaveFile(subPath string, file io.Reader) error {
	imgPath := path.Join(s.RootPath, subPath)
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

func (s *FileStorage) DeleteFile(subPath string) error {
	imgPath := path.Join(s.RootPath, subPath)
	if _, err := os.Stat(imgPath); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return os.Remove(imgPath)
}
