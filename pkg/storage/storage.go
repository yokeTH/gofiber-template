package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Config struct {
	BucketName      string `env:"BUCKET_NAME,required"`
	AccessKeyID     string `env:"ACCESS_KEY_ID,required"`
	AccessKeySecret string `env:"ACCESS_KEY_SECRET,required"`
	Endpoint        string `env:"ENDPOINT,required"`
	UrlPathStyle    bool   `env:"URL_PATH_STYLE" default:"false"`
}

type storage struct {
	presignClient *s3.PresignClient
	client        *s3.Client
	config        Config
	hostname      string
}

func appendSubdomain(baseURL, subdomain string) (string, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	hostParts := strings.Split(parsedURL.Host, ".")
	if len(hostParts) < 2 {
		return "", fmt.Errorf("invalid base URL: %s", baseURL)
	}

	newHost := fmt.Sprintf("%s.%s", subdomain, parsedURL.Host)
	parsedURL.Host = newHost

	return parsedURL.String(), nil
}

type Storage interface {
	UploadFile(ctx context.Context, key string, contentType string, file io.Reader) error
	GetSignedUrl(ctx context.Context, key string, expires time.Duration) (string, error)
	GetPublicUrl(key string) (string, error)
	DeleteFile(ctx context.Context, key string) error
}

func New(sConfig Config) (*storage, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(sConfig.AccessKeyID, sConfig.AccessKeySecret, "")),
		config.WithRegion("auto"),
		config.WithBaseEndpoint(sConfig.Endpoint),
	)

	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(client)

	newURL, err := appendSubdomain(sConfig.Endpoint, sConfig.BucketName)
	if err != nil {
		return nil, err
	}

	st := &storage{
		client:        client,
		presignClient: presignClient,
		config:        sConfig,
		hostname:      newURL,
	}

	return st, nil
}

func (s *storage) UploadFile(ctx context.Context, key string, contentType string, file io.Reader) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.config.BucketName),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
		Body:        file,
	})

	return err
}

func (s *storage) GetSignedUrl(ctx context.Context, key string, expires time.Duration) (string, error) {
	req, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.config.BucketName),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expires))

	if err != nil {
		return "", err
	}

	return req.URL, nil
}

func (s *storage) GetPublicUrl(key string) (string, error) {
	fmt.Println(s.config.UrlPathStyle)
	if s.config.UrlPathStyle {
		return url.JoinPath(s.config.Endpoint, s.config.BucketName, key)
	}
	return url.JoinPath(s.hostname, key)
}

func (s *storage) DeleteFile(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.config.BucketName),
		Key:    aws.String(key),
	})

	return err
}
