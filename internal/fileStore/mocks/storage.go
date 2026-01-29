package mock

import (
	"bytes"
)

type FileStorage struct{}

func (s *FileStorage) SaveFile(subPath string, file *bytes.Buffer) error {
	return nil
}

func (s *FileStorage) DeleteFile(key string) error {
	return nil
}

func (s *FileStorage) Type() string {
	return "Mock"
}
