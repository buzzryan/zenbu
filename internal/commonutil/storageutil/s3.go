package storageutil

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/buzzryan/zenbu/internal/config"
)

type s3Storage struct {
	client                   *s3.Client
	presignClient            *s3.PresignClient
	bucket                   string
	privateDir               string
	publicDir                string
	publicCloudfrontEndpoint string
}

func NewS3Storage(awsCfg aws.Config, cfg config.S3Config) Storage {
	client := s3.NewFromConfig(awsCfg)
	presignClient := s3.NewPresignClient(client)

	return &s3Storage{
		client:                   client,
		presignClient:            presignClient,
		bucket:                   cfg.Bucket,
		privateDir:               cfg.PrivateDir,
		publicDir:                cfg.PublicDir,
		publicCloudfrontEndpoint: cfg.PublicCloudfrontEndpoint,
	}
}

func (s *s3Storage) objectKey(scope Scope, filepath string) string {
	switch scope {
	case Private:
		return s.privateDir + "/" + filepath
	case Public:
		return s.publicDir + "/" + filepath
	default:
		return ""
	}
}

func (s *s3Storage) CreateUploadURL(ctx context.Context, scope Scope, filepath string) (url string, err error) {
	if filepath == "" {
		return "", errors.New("filepath required")
	}

	key := s.objectKey(scope, filepath)
	if key == "" {
		return "", errors.New("invalid scope")
	}

	res, err := s.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	}, s3.WithPresignExpires(time.Minute))
	if err != nil {
		return "", err
	}

	return res.URL, nil
}

func (s *s3Storage) GetPublicFileURL(ctx context.Context, filepath string) (string, error) {
	key := s.publicDir + "/" + filepath
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})

	if err != nil {
		return "", err
	}

	return s.publicCloudfrontEndpoint + "/" + filepath, nil
}
