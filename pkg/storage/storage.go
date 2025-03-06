package storage

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage struct {
	*s3.PresignClient
	*s3.Client
	BucketName string `env:"BUCKET_NAME"`
}

type IStorage interface {
	UploadFile(ctx context.Context, key string, contentType string, file io.Reader) error
	GetSignedUrl(ctx context.Context, key string, expires time.Duration) (string, error)
	GetPublicUrl(key string) string
	DeleteFile(ctx context.Context, key string) error
}
