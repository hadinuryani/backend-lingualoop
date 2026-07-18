package class

import "time"

type Class struct {
	ID                string
	AcademicYearID    string
	MajorID           *string
	LevelID           string
	GradeLevel        int
	ClassNumber       int
	ClassName         string
	Classroom         *string
	Capacity          int
	HomeroomTeacherID *string
	IsActive          bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
