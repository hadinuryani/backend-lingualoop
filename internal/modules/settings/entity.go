package settings

import "time"

type SchoolSettings struct {
	ID                        int
	SchoolName                string
	SchoolNPSN                string
	SchoolAddress             string
	SchoolPhone               string
	SchoolEmail               string
	SchoolLogoFileID          *string
	EducationLogoFileID       *string
	PrincipalName             string
	PrincipalNIP              string
	PrincipalSignatureFileID  *string
	MaxStudentsPerClass       int
	GradingSystem             string
	PassingGrade              int
	AppName                   string
	EnableStudentLogin        bool
	EnableTeacherLogin        bool
	MaintenanceMode           bool
	CreatedAt                 time.Time
	UpdatedAt                 time.Time

	// Relational data (joined from files table)
	SchoolLogoPath          *string
	EducationLogoPath       *string
	PrincipalSignaturePath  *string
}
