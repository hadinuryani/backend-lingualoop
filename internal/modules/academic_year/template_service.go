package academic_year

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type TemplateService interface {
	GenerateNext(ctx context.Context) (*TemplateResponse, error)
}

type templateService struct {
	repo Repository
}

func NewTemplateService(repo Repository) TemplateService {
	return &templateService{
		repo: repo,
	}
}

func (s *templateService) GenerateNext(ctx context.Context) (*TemplateResponse, error) {
	// 1. Cek apakah ada Draft yang menggantung
	draftExists, err := s.repo.CheckDraftExists(ctx)
	if err != nil {
		return nil, err
	}
	if draftExists {
		return nil, ErrDraftExists
	}

	// 2. Ambil tahun akademik terakhir (ORDER BY year DESC)
	latest, err := s.repo.GetLatestAcademicYearByDomain(ctx)
	if err != nil {
		if err == ErrAcademicYearNotFound {
			// Jika kosong, kembalikan template default
			return s.generateDefaultTemplate(), nil
		}
		return nil, err
	}

	// 3. Shift +1 Tahun
	template := TemplateData{
		Year:                     s.shiftYearString(latest.Year),
		StartDate:                s.shiftDate(latest.StartDate),
		EndDate:                  s.shiftDate(latest.EndDate),
		SemesterGanjilStart:      s.shiftNullDate(latest.SemGanjilStartDate),
		SemesterGanjilKbm:        s.shiftNullDate(latest.SemGanjilEndKBM),
		SemesterGanjilAssessment: s.shiftNullDate(latest.SemGanjilEndAssessment),
		SemesterGenapStart:       s.shiftNullDate(latest.SemGenapStartDate),
		SemesterGenapKbm:         s.shiftNullDate(latest.SemGenapEndKBM),
		SemesterGenapAssessment:  s.shiftNullDate(latest.SemGenapEndAssessment),
	}

	return &TemplateResponse{
		SourceYear: latest.Year,
		Template:   template,
	}, nil
}

func (s *templateService) shiftDate(d time.Time) string {
	if d.IsZero() {
		return ""
	}
	// Tambah tepat 1 tahun
	shifted := d.AddDate(1, 0, 0)
	return shifted.Format("2006-01-02")
}

func (s *templateService) shiftNullDate(nd sql.NullTime) string {
	if !nd.Valid {
		return ""
	}
	shifted := nd.Time.AddDate(1, 0, 0)
	return shifted.Format("2006-01-02")
}

func (s *templateService) shiftYearString(year string) string {
	// Contoh format: "2025/2026"
	parts := strings.Split(year, "/")
	if len(parts) != 2 {
		return year // fallback
	}

	y1, err1 := strconv.Atoi(parts[0])
	y2, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return year
	}

	return fmt.Sprintf("%d/%d", y1+1, y2+1)
}

func (s *templateService) generateDefaultTemplate() *TemplateResponse {
	now := time.Now()
	yearStr := fmt.Sprintf("%d/%d", now.Year(), now.Year()+1)
	
	return &TemplateResponse{
		SourceYear: "Sistem (Default)",
		Template: TemplateData{
			Year: yearStr,
		},
	}
}
