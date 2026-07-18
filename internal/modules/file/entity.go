package file

import (
	"time"
)

type File struct {
	ID           string
	ResourceType string
	StoragePath  string
	OriginalName string
	MimeType     string
	SizeBytes    int64
	UploadedBy   *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
