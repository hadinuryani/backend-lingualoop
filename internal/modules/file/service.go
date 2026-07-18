package file

import (
	"context"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"backend-lingualoop/pkg/storage"
	"github.com/google/uuid"
	"log/slog"
)

type Service interface {
	UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader, resourceType string, uploadedBy *string) (*FileResponse, error)
	GetFile(ctx context.Context, id string) (*FileResponse, error)
}

type service struct {
	repo    Repository
	storage storage.Storage
}

func NewService(repo Repository, store storage.Storage) Service {
	return &service{
		repo:    repo,
		storage: store,
	}
}

// 2MB limit
const MaxImageSize = 2 * 1024 * 1024

func (s *service) UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader, resourceType string, uploadedBy *string) (*FileResponse, error) {
	if header.Size > MaxImageSize {
		return nil, ErrFileTooLarge
	}

	// Read first 512 bytes for MIME detection
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, ErrSystemFail
	}

	mimeType := http.DetectContentType(buffer)
	if mimeType != "image/jpeg" && mimeType != "image/png" && mimeType != "image/webp" {
		return nil, ErrInvalidFileType
	}

	// Reset file pointer to beginning for DecodeConfig
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, ErrSystemFail
	}

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		slog.Warn("Failed to decode image config (possible malicious file)", "error", err)
		return nil, ErrInvalidFileType
	}

	if config.Width < 200 || config.Height < 200 {
		return nil, ErrInvalidFileType // Or custom ErrImageTooSmall
	}
	if config.Width > 4000 || config.Height > 4000 {
		return nil, ErrFileTooLarge // Or custom ErrImageTooLarge
	}

	// Reset pointer again for saving
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, ErrSystemFail
	}

	folder := "public/" + resourceType // e.g. "public/majors"

	// Save via storage interface
	info, err := s.storage.Save(ctx, file, folder, header.Filename, mimeType)
	if err != nil {
		slog.Error("Failed to save file to storage", "error", err)
		return nil, ErrSystemFail
	}

	now := time.Now()
	f := &File{
		ID:           uuid.New().String(),
		ResourceType: resourceType,
		StoragePath:  info.Path,
		OriginalName: header.Filename,
		MimeType:     info.MimeType,
		SizeBytes:    info.Size,
		UploadedBy:   uploadedBy,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.Create(ctx, f); err != nil {
		slog.Error("Failed to save file to DB", "error", err)
		// Rollback storage if DB fails
		_ = s.storage.Delete(ctx, info.Path)
		return nil, ErrSystemFail
	}

	return &FileResponse{
		ID:           f.ID,
		URL:          s.storage.GetURL(f.StoragePath),
		OriginalName: f.OriginalName,
		MimeType:     f.MimeType,
		SizeBytes:    f.SizeBytes,
		CreatedAt:    f.CreatedAt,
	}, nil
}

func (s *service) GetFile(ctx context.Context, id string) (*FileResponse, error) {
	f, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &FileResponse{
		ID:           f.ID,
		URL:          s.storage.GetURL(f.StoragePath),
		OriginalName: f.OriginalName,
		MimeType:     f.MimeType,
		SizeBytes:    f.SizeBytes,
		CreatedAt:    f.CreatedAt,
	}, nil
}
