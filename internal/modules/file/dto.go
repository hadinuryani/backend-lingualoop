package file

import "time"

type FileResponse struct {
	ID           string    `json:"id"`
	URL          string    `json:"url"` // Full absolute URL to access the file
	OriginalName string    `json:"original_name"`
	MimeType     string    `json:"mime_type"`
	SizeBytes    int64     `json:"size_bytes"`
	CreatedAt    time.Time `json:"created_at"`
}
