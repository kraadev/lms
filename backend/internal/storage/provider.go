package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

// Provider abstracts storage backends (Local Disk, S3/MinIO)
type Provider interface {
	Save(ctx context.Context, key string, r io.Reader) (string, error)
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) bool
}

// LocalDiskProvider implements Provider for local filesystem storage
type LocalDiskProvider struct {
	baseDir string
}

func NewLocalDiskProvider(baseDir string) (*LocalDiskProvider, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}
	return &LocalDiskProvider{baseDir: baseDir}, nil
}

func (p *LocalDiskProvider) Save(ctx context.Context, key string, r io.Reader) (string, error) {
	destPath := filepath.Join(p.baseDir, filepath.Clean(key))
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return "", err
	}

	out, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, r); err != nil {
		return "", err
	}

	return key, nil
}

func (p *LocalDiskProvider) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	destPath := filepath.Join(p.baseDir, filepath.Clean(key))
	return os.Open(destPath)
}

func (p *LocalDiskProvider) Delete(ctx context.Context, key string) error {
	destPath := filepath.Join(p.baseDir, filepath.Clean(key))
	return os.Remove(destPath)
}

func (p *LocalDiskProvider) Exists(ctx context.Context, key string) bool {
	destPath := filepath.Join(p.baseDir, filepath.Clean(key))
	info, err := os.Stat(destPath)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
