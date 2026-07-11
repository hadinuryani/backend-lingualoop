package academic_year

import "time"

type AcademicYearRequest struct {
	Year                     string `json:"year" binding:"required"`
	StartDate                string `json:"start_date" binding:"required"`
	EndDate                  string `json:"end_date" binding:"required"`
	SemesterGanjilStart      string `json:"semester_ganjil_start"`
	SemesterGanjilKbm        string `json:"semester_ganjil_kbm"`
	SemesterGanjilAssessment string `json:"semester_ganjil_assessment"`
	SemesterGenapStart       string `json:"semester_genap_start"`
	SemesterGenapKbm         string `json:"semester_genap_kbm"`
	SemesterGenapAssessment  string `json:"semester_genap_assessment"`
}

type SemesterStatusRequest struct {
	Semester string `json:"semester" binding:"required,oneof=ganjil genap"` // ganjil atau genap
	Status   string `json:"status" binding:"required"`
}

type CloseSemesterRequest struct {
	Semester string `json:"semester" binding:"required,oneof=ganjil genap"` // ganjil atau genap
}

type StudentPromotion struct {
	StudentID string `json:"student_id" binding:"required"`
	ClassID   string `json:"class_id" binding:"required"` // Kelas saat ini
	Status    string `json:"status" binding:"required,oneof=Lulus 'Naik Kelas' 'Tinggal Kelas' 'Tidak Lulus'"`
}

type FinalizePromotionRequest struct {
	TargetYearID string             `json:"target_year_id" binding:"required"`
	Promotions   []StudentPromotion `json:"promotions" binding:"required"`
}

// -------------------------------------------------------------
// RESPONSES
// -------------------------------------------------------------

type SemesterData struct {
	StartDate     string `json:"start_date,omitempty"`
	EndKBM        string `json:"end_kbm,omitempty"`
	EndAssessment string `json:"end_assessment,omitempty"`
	Status        string `json:"status"`
}

type AcademicYearResponse struct {
	ID             string       `json:"id"`
	Year           string       `json:"year"`
	StartDate      string       `json:"start_date"`
	EndDate        string       `json:"end_date"`
	Status         string       `json:"status"`
	IsActive       bool         `json:"is_active"`
	Semester       string       `json:"semester"` // "GANJIL", "GENAP", atau "-"
	SemesterGanjil SemesterData `json:"semester_ganjil"`
	SemesterGenap  SemesterData `json:"semester_genap"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}
