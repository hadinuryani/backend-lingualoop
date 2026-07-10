package class

import "time"

type ClassRequest struct {
	AcademicYearID    string `json:"academic_year_id" binding:"required"`
	MajorID           string `json:"major_id" binding:"required"`
	LevelID           string `json:"level_id" binding:"required"`
	ClassName         string `json:"class_name" binding:"required"`
	Classroom         string `json:"classroom"`
	Capacity          int    `json:"capacity"`
	HomeroomTeacherID string `json:"homeroom_teacher_id"`
}

type ClassBatchRequest struct {
	AcademicYearID string   `json:"academic_year_id" binding:"required"`
	MajorID        string   `json:"major_id" binding:"required"`
	MajorCode      string   `json:"major_code" binding:"required"` // Untuk digabung dalam nama, ex: IPA
	LevelIDs       []string `json:"level_ids" binding:"required,min=1"`
	Capacity       int      `json:"capacity"`
	Count          int      `json:"count" binding:"required,min=1"`
}

type ClassUpdateRequest struct {
	Capacity          int    `json:"capacity"`
	HomeroomTeacherID string `json:"homeroom_teacher_id"`
}

type ClassResponse struct {
	ID                string    `json:"id"`
	AcademicYearID    string    `json:"academic_year_id"`
	MajorID           string    `json:"major_id,omitempty"`
	LevelID           string    `json:"level_id"`
	ClassName         string    `json:"class_name"`
	Classroom         string    `json:"classroom,omitempty"`
	Capacity          int       `json:"capacity"`
	HomeroomTeacherID string    `json:"homeroom_teacher_id,omitempty"`
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
