package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	Client *s3.Client
	Bucket string
}

func NewS3Storage(bucket string) (*S3Storage, error) {

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
	)

	if err != nil {
		return nil, err
	}

	return &S3Storage{
		Client: s3.NewFromConfig(cfg),
		Bucket: bucket,
	}, nil
}

func (s *S3Storage) Upload(
	key string,
	body io.Reader,
) (string, error) {

	_, err := s.Client.PutObject(
		context.TODO(),
		&s3.PutObjectInput{
			Bucket: &s.Bucket,
			Key:    &key,
			Body:   body,
		},
	)

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf(
		"https://%s.s3.amazonaws.com/%s",
		s.Bucket,
		key,
	)

	return url, nil
}