package settings

import "time"

type SchoolSettingsRequest struct {
	SchoolName                string  `json:"school_name" binding:"required"`
	SchoolNPSN                string  `json:"school_npsn"`
	SchoolAddress             string  `json:"school_address"`
	SchoolPhone               string  `json:"school_phone"`
	SchoolEmail               string  `json:"school_email" binding:"omitempty,email"`
	SchoolLogoFileID          *string `json:"school_logo_file_id"`
	EducationLogoFileID       *string `json:"education_logo_file_id"`
	PrincipalName             string  `json:"principal_name"`
	PrincipalNIP              string  `json:"principal_nip"`
	PrincipalSignatureFileID  *string `json:"principal_signature_file_id"`
	MaxStudentsPerClass       int     `json:"max_students_per_class" binding:"min=1"`
	GradingSystem             string  `json:"grading_system" binding:"oneof=numeric letter"`
	PassingGrade              int     `json:"passing_grade" binding:"min=0,max=100"`
	AppName                   string  `json:"app_name" binding:"required"`
	EnableStudentLogin        bool    `json:"enable_student_login"`
	EnableTeacherLogin        bool    `json:"enable_teacher_login"`
	MaintenanceMode           bool    `json:"maintenance_mode"`
}

type SchoolSettingsResponse struct {
	ID                        int       `json:"id"`
	SchoolName                string    `json:"school_name"`
	SchoolNPSN                string    `json:"school_npsn"`
	SchoolAddress             string    `json:"school_address"`
	SchoolPhone               string    `json:"school_phone"`
	SchoolEmail               string    `json:"school_email"`
	SchoolLogoFileID          *string   `json:"school_logo_file_id,omitempty"`
	SchoolLogoURL             string    `json:"school_logo_url,omitempty"`
	EducationLogoFileID       *string   `json:"education_logo_file_id,omitempty"`
	EducationLogoURL          string    `json:"education_logo_url,omitempty"`
	PrincipalName             string    `json:"principal_name"`
	PrincipalNIP              string    `json:"principal_nip"`
	PrincipalSignatureFileID  *string   `json:"principal_signature_file_id,omitempty"`
	PrincipalSignatureURL     string    `json:"principal_signature_url,omitempty"`
	MaxStudentsPerClass       int       `json:"max_students_per_class"`
	GradingSystem             string    `json:"grading_system"`
	PassingGrade              int       `json:"passing_grade"`
	AppName                   string    `json:"app_name"`
	EnableStudentLogin        bool      `json:"enable_student_login"`
	EnableTeacherLogin        bool      `json:"enable_teacher_login"`
	MaintenanceMode           bool      `json:"maintenance_mode"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}
