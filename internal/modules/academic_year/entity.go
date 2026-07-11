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

	SemesterOdd  = "ganjil"
	SemesterEven = "genap"
)
