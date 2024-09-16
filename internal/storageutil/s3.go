package storageutil

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Storage struct {
	client                   *s3.Client
	presignClient            *s3.PresignClient
	bucket                   string
	privateDir               string
	publicDir                string
	publicCloudfrontEndpoint string
}

func NewS3Storage(client *s3.Client, presignClient *s3.PresignClient, bucket, privateDir, publicDir, publicCloudfrontEndpoint string) Storage {
	return &s3Storage{
		client:                   client,
		presignClient:            presignClient,
		bucket:                   bucket,
		privateDir:               privateDir,
		publicDir:                publicDir,
		publicCloudfrontEndpoint: publicCloudfrontEndpoint,
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

func (s *s3Storage) GetUploadURL(ctx context.Context, scope Scope, filepath string) (url string, err error) {
	if filepath == "" {
		return "", errors.New("filepath required")
	}

	key := s.objectKey(scope, filepath)
	if key == "" {
		return "", errors.New("invalid scope")
	}

	expiresAt := time.Now().Add(1 * time.Minute)

	res, err := s.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:  &s.bucket,
		Key:     &key,
		Expires: &expiresAt,
	})
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

	return s.publicCloudfrontEndpoint + "/" + key, nil
}
