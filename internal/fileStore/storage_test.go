package fileStore

import (
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func TestUploadUrl(t *testing.T) {
	godotenv.Load("../../.env")

	endpoint := os.Getenv("DEV_S3_ENDPOINT")
	key := os.Getenv("DEV_S3_KEY")
	secret := os.Getenv("DEV_S3_SECRET")
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String("us-east-1"),
		S3ForcePathStyle: aws.Bool(false),
	}

	s3Client := s3.New(session.Must(session.NewSession(s3Config)))
	bucket := os.Getenv("DEV_S3_BUCKET")
	fileStore := S3Storage{
		Client:     s3Client,
		BucketName: bucket,
	}

	fname := uuid.New().String() + ".jpg"

	const MAX_FILE_SIZE = 262_144_000 // 250 MiB
	url, err := fileStore.PresignedUploadURL(fname, 30*time.Minute, MAX_FILE_SIZE)
	if err != nil {
		t.Error(err)
	}

	t.Logf("UPLOAD URL: %s", url)
}
