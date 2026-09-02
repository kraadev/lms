package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"lms/internal/utils"
)

type StorageService struct {
	baseDir string
}

func NewStorageService(baseDir string) (*StorageService, error) {
	categories := []string{"materials", "assignments", "submissions", "avatars"}
	for _, c := range categories {
		dir := filepath.Join(baseDir, c)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create upload directory %s: %w", dir, err)
		}
	}
	return &StorageService{baseDir: baseDir}, nil
}

func (s *StorageService) SaveUploadedFile(fileHeader *multipart.FileHeader, category string) (string, error) {
	allowedCategories := map[string]bool{
		"materials":   true,
		"assignments": true,
		"submissions": true,
		"avatars":     true,
	}
	if !allowedCategories[category] {
		return "", fmt.Errorf("invalid upload category: %s", category)
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedExtensions := map[string]bool{
		".pdf": true, ".doc": true, ".docx": true, ".ppt": true, ".pptx": true,
		".xls": true, ".xlsx": true, ".zip": true, ".tar": true, ".gz": true,
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
		".mp4": true, ".mp3": true, ".txt": true, ".md": true, ".go": true,
		".py": true, ".js": true, ".ts": true, ".json": true,
	}

	if ext != "" && !allowedExtensions[ext] {
		return "", fmt.Errorf("file extension '%s' is not allowed", ext)
	}

	uniqueName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	targetDir := filepath.Join(s.baseDir, category)
	targetPath := filepath.Join(targetDir, uniqueName)

	cleanPath := filepath.Clean(targetPath)
	if !strings.HasPrefix(cleanPath, filepath.Clean(s.baseDir)) {
		return "", fmt.Errorf("illegal file path traversal attempt")
	}

	dst, err := os.Create(cleanPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	// Returns relative path like "uploads/materials/uuid.pdf"
	return filepath.ToSlash(filepath.Join(s.baseDir, category, uniqueName)), nil
}

func (s *StorageService) ServeFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// e.g. /api/files/materials/uuid.pdf
		relPath := strings.TrimPrefix(r.URL.Path, "/api/files/")
		cleanPath := filepath.Clean(filepath.Join(s.baseDir, relPath))

		if !strings.HasPrefix(cleanPath, filepath.Clean(s.baseDir)) {
			utils.Forbidden(w, "Access denied")
			return
		}

		if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
			utils.NotFound(w, "File not found")
			return
		}

		http.ServeFile(w, r, cleanPath)
	}
}
