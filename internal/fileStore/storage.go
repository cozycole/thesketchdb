package fileStore

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

type FileStorageInterface interface {
	DeleteFile(string) error
	Exists(string) (bool, error)
	PresignedUploadURL(string, time.Duration, int) (string, error)
	SaveFile(string, *bytes.Buffer) error
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

func (s *S3Storage) Exists(s3Key string) (bool, error) {
	_, err := s.Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(s3Key),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				// The object does not exist
				return false, nil
			default:
				// Other AWS error (e.g., permissions, network issues)
				return false, fmt.Errorf("AWS error: %w", err)
			}
		}
		// Non-AWS error
		return false, fmt.Errorf("non-AWS error: %w", err)
	}

	return true, nil
}

func (s *S3Storage) PresignedUploadURL(filename string, duration time.Duration, contentLength int) (string, error) {
	putObject := &s3.PutObjectInput{
		Bucket:        aws.String(s.BucketName),
		Key:           aws.String(filename),
		ContentLength: aws.Int64(int64(contentLength)),
	}

	req, _ := s.Client.PutObjectRequest(putObject)
	url, err := req.Presign(duration)
	if err != nil {
		return "", err
	}

	return url, nil
}
