package img

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type FileStorageInterface interface {
	SaveFile(string, *bytes.Buffer) error
	DeleteFile(string) error
	Type() string
}

type FileStorage struct {
	RootPath string
}

func (s *FileStorage) SaveFile(subPath string, file *bytes.Buffer) error {
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

	body := bytes.NewReader(file.Bytes())
	if _, err := io.Copy(dst, body); err != nil {
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

func (s *FileStorage) Type() string {
	return "Local"
}

type S3Storage struct {
	Client     *s3.S3
	BucketName string
}

func (s *S3Storage) SaveFile(subPath string, file *bytes.Buffer) error {
	body := bytes.NewReader(file.Bytes())
	_, err := s.Client.PutObject(&s3.PutObjectInput{
		Bucket:       &s.BucketName,
		Key:          &subPath,
		Body:         body,
		ACL:          aws.String("public-read"),
		ContentType:  aws.String(http.DetectContentType(file.Bytes())),
		CacheControl: aws.String("public, max-age=31536000, immutable"),
	})
	return err
}

func (s *S3Storage) DeleteFile(subPath string) error {
	_, err := s.Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &s.BucketName,
		Key:    &subPath,
	})
	return err
}

func (s *S3Storage) Type() string {
	return "S3"
}
