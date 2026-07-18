package major

import "time"

type MajorRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	LogoFileID  string `json:"logo_file_id"` // Link to files.id
}

type MajorResponse struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	LogoURL     string    `json:"logo_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
