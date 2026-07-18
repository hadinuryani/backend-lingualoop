package academic_year

import (
	"database/sql"
	"time"
)

type AcademicYear struct {
	ID                     string
	Year                   string
	StartDate              time.Time
	EndDate                time.Time
	Status                 string // Draft, Aktif, Menunggu Kenaikan, Selesai
	SemGanjilStartDate     sql.NullTime
	SemGanjilEndKBM        sql.NullTime
	SemGanjilEndAssessment sql.NullTime
	SemGanjilStatus        string // Belum Aktif, Aktif, Masa Penilaian, Siap Ditutup, Terkunci
	SemGenapStartDate      sql.NullTime
	SemGenapEndKBM         sql.NullTime
	SemGenapEndAssessment  sql.NullTime
	SemGenapStatus         string
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

const (
	StatusDraft            = "Draft"
	StatusActive           = "Aktif"
	StatusPendingPromotion = "Menunggu Kenaikan"
	StatusFinished         = "Selesai"

	SemStatusNotActive    = "Belum Aktif"
	SemStatusActive       = "Aktif"
	SemStatusAssessment   = "Masa Penilaian"
	SemStatusReadyToClose = "Siap Ditutup"
	SemStatusLocked       = "Terkunci"

	SemesterOddKey  = "ganjil"
	SemesterEvenKey = "genap"

	SemesterOddLabel  = "GANJIL"
	SemesterEvenLabel = "GENAP"

	JobStatusPending = "PENDING"
	JobStatusRunning = "RUNNING"
	JobStatusDone    = "DONE"
	JobStatusFailed  = "FAILED"

	PromotionStatusPromoted  = "PROMOTED"
	PromotionStatusRetained  = "RETAINED"
	PromotionStatusGraduated = "GRADUATED"
	PromotionStatusFailed    = "FAILED"
)

type PromotionJob struct {
	ID              string
	AcademicYearID  string
	Status          string
	TotalStudents   int
	SuccessStudents int
	FailedStudents  int
	ExecutedBy      *string
	StartedAt       *time.Time
	FinishedAt      *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type StudentPromotionHistory struct {
	ID                 string
	StudentID          string
	FromClassID        *string
	ToClassID          *string
	FromAcademicYearID string
	ToAcademicYearID   *string
	Status             string
	CreatedAt          time.Time
}
