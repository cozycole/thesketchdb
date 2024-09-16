package img

import (
	"io"
	"mime/multipart"
	"os"
	"path"
)

type FileStorageInterface interface {
	SaveMultipartFile(string, multipart.File) error
}

type FileStorage struct {
	Path string
}

func (s *FileStorage) SaveMultipartFile(subPath string, file multipart.File) error {
	imgPath := path.Join(s.Path, subPath)
	dst, err := os.Create(imgPath)
	if err != nil {
		return err
	}

	if _, err := io.Copy(dst, file); err != nil {
		return err
	}
	return nil
}
