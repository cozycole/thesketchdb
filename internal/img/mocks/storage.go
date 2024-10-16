package mock

import (
	"io"
)

type FileStorage struct{}

func (s *FileStorage) SaveFile(subPath string, file io.Reader) error {
	return nil
}
