package mock

import (
	"mime/multipart"
)

type FileStorage struct{}

func (s *FileStorage) SaveMultipartFile(subPath string, file multipart.File) error {
	return nil
}
