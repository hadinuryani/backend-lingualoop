package subject

import "time"

type SubjectRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	MajorID     string `json:"major_id"` // Boleh kosong (umum)
	LevelID     string `json:"level_id"` // Boleh kosong
}

type SubjectResponse struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	MajorID     string    `json:"major_id,omitempty"`
	LevelID     string    `json:"level_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
