package storage

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"
)

func TestLocalDiskProvider(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lms-storage-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	provider, err := NewLocalDiskProvider(tmpDir)
	if err != nil {
		t.Fatalf("failed to init local provider: %v", err)
	}

	ctx := context.Background()
	key := "assignments/submission1.txt"
	data := "LMS student homework content"

	// 1. Save
	savedKey, err := provider.Save(ctx, key, bytes.NewReader([]byte(data)))
	if err != nil {
		t.Fatalf("failed to save file: %v", err)
	}
	if savedKey != key {
		t.Errorf("expected saved key '%s', got '%s'", key, savedKey)
	}

	// 2. Exists
	if !provider.Exists(ctx, key) {
		t.Errorf("expected file to exist")
	}

	// 3. Get
	reader, err := provider.Get(ctx, key)
	if err != nil {
		t.Fatalf("failed to get file: %v", err)
	}
	readBytes, _ := io.ReadAll(reader)
	reader.Close()
	if string(readBytes) != data {
		t.Errorf("expected '%s', got '%s'", data, string(readBytes))
	}

	// 4. Delete
	if err := provider.Delete(ctx, key); err != nil {
		t.Fatalf("failed to delete file: %v", err)
	}
	if provider.Exists(ctx, key) {
		t.Errorf("expected file to be deleted")
	}
}
