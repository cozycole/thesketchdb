package img

import (
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
