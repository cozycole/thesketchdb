package mock

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path"
)

type FileStorage struct{}

func (s *FileStorage) SaveFile(subPath string, file *bytes.Buffer) error {
	imgPath := path.Join(os.TempDir(), subPath)
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

func (s *FileStorage) DeleteFile(key string) error {
	return nil
}

func (s *FileStorage) Type() string {
	return "Mock"
}
