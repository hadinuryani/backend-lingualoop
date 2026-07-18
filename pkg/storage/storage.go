package storage

import (
	"context"
	"io"
)

// Info represents metadata of a saved file
type Info struct {
	Path        string // e.g., "public/majors/uuid.png"
	Size        int64
	MimeType    string
	OriginalExt string
}

// Storage is the core interface for all file operations
type Storage interface {
	// Save stores a file and returns its metadata
	Save(ctx context.Context, reader io.Reader, folder, originalFilename, mimeType string) (Info, error)
	
	// Delete removes a file
	Delete(ctx context.Context, path string) error
	
	// GetURL converts a storage path to a full accessible URL
	GetURL(path string) string
}
