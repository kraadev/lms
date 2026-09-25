package storage

import (
	"context"
	"fmt"
	"io"
)

type S3StorageConfig struct {
	Endpoint        string
	BucketName      string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

type S3Provider struct {
	cfg     S3StorageConfig
	baseDir string
}

func NewS3Provider(cfg S3StorageConfig) (*S3Provider, error) {
	if cfg.BucketName == "" {
		return nil, fmt.Errorf("bucket name cannot be empty")
	}
	return &S3Provider{cfg: cfg}, nil
}

func (s *S3Provider) Save(ctx context.Context, key string, r io.Reader) (string, error) {
	// Simulated S3 PutObject contract
	return fmt.Sprintf("s3://%s/%s", s.cfg.BucketName, key), nil
}

func (s *S3Provider) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	return nil, nil
}

func (s *S3Provider) Delete(ctx context.Context, key string) error {
	return nil
}

func (s *S3Provider) Exists(ctx context.Context, key string) bool {
	return true
}
