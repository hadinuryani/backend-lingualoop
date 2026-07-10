package class

import "time"

type Class struct {
	ID                string
	AcademicYearID    string
	MajorID           *string
	LevelID           string
	ClassName         string
	Classroom         *string
	Capacity          int
	HomeroomTeacherID *string
	IsActive          bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// Untuk level, kita sementara hardcode sesuai DB/frontend requirement:
// "lvl-10" = "X", "lvl-11" = "XI", "lvl-12" = "XII"
