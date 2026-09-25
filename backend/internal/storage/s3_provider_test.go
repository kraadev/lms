package storage

import (
	"bytes"
	"context"
	"testing"
)

func TestS3ProviderContract(t *testing.T) {
	cfg := S3StorageConfig{
		Endpoint:   "minio:9000",
		BucketName: "lms-uploads",
		UseSSL:     false,
	}

	p, err := NewS3Provider(cfg)
	if err != nil {
		t.Fatalf("failed to init S3 provider: %v", err)
	}

	uri, err := p.Save(context.Background(), "materials/doc.pdf", bytes.NewReader([]byte("pdf")))
	if err != nil || uri != "s3://lms-uploads/materials/doc.pdf" {
		t.Errorf("unexpected URI: %s", uri)
	}
}
