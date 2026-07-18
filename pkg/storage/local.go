package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type LocalStorage struct {
	BaseDir string
	BaseURL string
}

// NewLocalStorage initializes the local disk storage
func NewLocalStorage(baseDir, baseURL string) Storage {
	return &LocalStorage{
		BaseDir: baseDir,
		BaseURL: baseURL,
	}
}

func (l *LocalStorage) Save(ctx context.Context, reader io.Reader, folder, originalFilename, mimeType string) (Info, error) {
	// folder e.g., "public/majors"
	ext := strings.ToLower(filepath.Ext(originalFilename))
	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// e.g. "public/majors/uuid.png"
	relativePath := filepath.Join(folder, newFilename) 
	
	// e.g. "./uploads/public/majors/uuid.png"
	fullPath := filepath.Join(l.BaseDir, relativePath)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return Info{}, fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	file, err := os.Create(fullPath)
	if err != nil {
		return Info{}, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write data securely with context cancellation support
	buf := make([]byte, 32*1024)
	var size int64
	for {
		select {
		case <-ctx.Done():
			file.Close()
			os.Remove(fullPath)
			return Info{}, ctx.Err()
		default:
		}

		nr, er := reader.Read(buf)
		if nr > 0 {
			nw, ew := file.Write(buf[0:nr])
			if nw > 0 {
				size += int64(nw)
			}
			if ew != nil {
				file.Close()
				os.Remove(fullPath)
				return Info{}, fmt.Errorf("failed to write file content: %w", ew)
			}
		}
		if er != nil {
			if er != io.EOF {
				file.Close()
				os.Remove(fullPath)
				return Info{}, fmt.Errorf("failed to read from source: %w", er)
			}
			break
		}
	}

	// Always use forward slashes for relative path in DB
	dbPath := filepath.ToSlash(relativePath)

	return Info{
		Path:        dbPath,
		Size:        size,
		MimeType:    mimeType,
		OriginalExt: ext,
	}, nil
}

func (l *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(l.BaseDir, path)
	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil // Already deleted
		}
		return err
	}
	return nil
}

func (l *LocalStorage) GetURL(path string) string {
	if path == "" {
		return ""
	}
	// e.g., BaseURL = "http://localhost:8080/uploads"
	// path = "public/majors/uuid.png"
	// output = "http://localhost:8080/uploads/public/majors/uuid.png"
	
	baseURL := strings.TrimRight(l.BaseURL, "/")
	path = strings.TrimLeft(path, "/")
	
	return fmt.Sprintf("%s/%s", baseURL, path)
}
