package academic_year

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	GetAll(ctx context.Context) ([]*AcademicYearResponse, error)
	GetByID(ctx context.Context, id string) (*AcademicYearResponse, error)
	Create(ctx context.Context, req AcademicYearRequest) (*AcademicYearResponse, error)
	Update(ctx context.Context, id string, req AcademicYearRequest) (*AcademicYearResponse, error)
	Activate(ctx context.Context, id string) (*AcademicYearResponse, error)
	UpdateSemesterStatus(ctx context.Context, id string, req SemesterStatusRequest) (*AcademicYearResponse, error)
	CloseSemester(ctx context.Context, id string, req CloseSemesterRequest) (*AcademicYearResponse, error)
	Delete(ctx context.Context, id string) error

	// FinalizePromotion akan diimplementasikan nanti saat student_classes module tersedia
	// FinalizePromotion(ctx context.Context, id string, req FinalizePromotionRequest) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func parseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, errors.New("date is empty")
	}
	return time.Parse("2006-01-02", dateStr)
}

func parseNullDate(dateStr string) (sql.NullTime, error) {
	if dateStr == "" {
		return sql.NullTime{Valid: false}, nil
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return sql.NullTime{Valid: false}, ErrInvalidDateFormat
	}
	return sql.NullTime{Time: t, Valid: true}, nil
}

func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

func formatNullDate(nt sql.NullTime) string {
	if !nt.Valid {
		return ""
	}
	return nt.Time.Format("2006-01-02")
}

func (s *service) GetAll(ctx context.Context) ([]*AcademicYearResponse, error) {
	years, err := s.repo.FindAll(ctx)
	if err != nil {
		slog.Error("Failed to query academic years", "error", err)
		return nil, ErrSystemFail
	}

	var responses []*AcademicYearResponse
	for _, y := range years {
		responses = append(responses, mapEntityToDTO(y))
	}
	if responses == nil {
		responses = []*AcademicYearResponse{}
	}
	return responses, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*AcademicYearResponse, error) {
	y, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrAcademicYearNotFound) {
			return nil, ErrAcademicYearNotFound
		}
		slog.Error("Failed to query academic year by ID", "error", err, "id", id)
		return nil, ErrSystemFail
	}
	return mapEntityToDTO(y), nil
}

func (s *service) Create(ctx context.Context, req AcademicYearRequest) (*AcademicYearResponse, error) {
	existingID, err := s.repo.GetIDByYear(ctx, req.Year)
	if err != nil && !errors.Is(err, ErrAcademicYearNotFound) {
		slog.Error("Failed to check unique year", "error", err)
		return nil, ErrSystemFail
	}
	if existingID != "" {
		return nil, ErrAcademicYearExists
	}

	startDate, err := parseDate(req.StartDate)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}
	endDate, err := parseDate(req.EndDate)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	if endDate.Before(startDate) {
		return nil, ErrInvalidDate
	}

	ganjilStart, err := parseNullDate(req.SemesterGanjilStart)
	if err != nil {
		return nil, err
	}
	ganjilEndKBM, err := parseNullDate(req.SemesterGanjilKbm)
	if err != nil {
		return nil, err
	}
	ganjilEndAssesment, err := parseNullDate(req.SemesterGanjilAssessment)
	if err != nil {
		return nil, err
	}

	genapStart, err := parseNullDate(req.SemesterGenapStart)
	if err != nil {
		return nil, err
	}
	genapEndKBM, err := parseNullDate(req.SemesterGenapKbm)
	if err != nil {
		return nil, err
	}
	genapEndAssessment, err := parseNullDate(req.SemesterGenapAssessment)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	ay := &AcademicYear{
		ID:                     uuid.New().String(),
		Year:                   req.Year,
		StartDate:              startDate,
		EndDate:                endDate,
		Status:                 StatusDraft,
		SemGanjilStartDate:     ganjilStart,
		SemGanjilEndKBM:        ganjilEndKBM,
		SemGanjilEndAssessment: ganjilEndAssesment,
		SemGanjilStatus:        SemStatusNotActive,
		SemGenapStartDate:      genapStart,
		SemGenapEndKBM:         genapEndKBM,
		SemGenapEndAssessment:  genapEndAssessment,
		SemGenapStatus:         SemStatusNotActive,
		CreatedAt:              now,
		UpdatedAt:              now,
	}

	if err := s.repo.Create(ctx, ay); err != nil {
		slog.Error("Failed to create academic year", "error", err)
		return nil, ErrSystemFail
	}

	return mapEntityToDTO(ay), nil
}

func (s *service) Update(ctx context.Context, id string, req AcademicYearRequest) (*AcademicYearResponse, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrAcademicYearNotFound) {
			return nil, ErrAcademicYearNotFound
		}
		slog.Error("Failed to get year for update", "error", err, "id", id)
		return nil, ErrSystemFail
	}

	if req.Year != existing.Year {
		existingID, err := s.repo.GetIDByYear(ctx, req.Year)
		if err != nil && !errors.Is(err, ErrAcademicYearNotFound) {
			slog.Error("Failed to check unique year on update", "error", err)
			return nil, ErrSystemFail
		}
		if existingID != "" && existingID != id {
			return nil, ErrAcademicYearExists
		}
	}

	startDate, err := parseDate(req.StartDate)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}
	endDate, err := parseDate(req.EndDate)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	if endDate.Before(startDate) {
		return nil, ErrInvalidDate
	}

	ganjilStart, err := parseNullDate(req.SemesterGanjilStart)
	if err != nil {
		return nil, err
	}
	ganjilEndKBM, err := parseNullDate(req.SemesterGanjilKbm)
	if err != nil {
		return nil, err
	}
	ganjilEndAssesment, err := parseNullDate(req.SemesterGanjilAssessment)
	if err != nil {
		return nil, err
	}

	genapStart, err := parseNullDate(req.SemesterGenapStart)
	if err != nil {
		return nil, err
	}
	genapEndKBM, err := parseNullDate(req.SemesterGenapKbm)
	if err != nil {
		return nil, err
	}
	genapEndAssessment, err := parseNullDate(req.SemesterGenapAssessment)
	if err != nil {
		return nil, err
	}

	existing.Year = req.Year
	existing.StartDate = startDate
	existing.EndDate = endDate
	existing.SemGanjilStartDate = ganjilStart
	existing.SemGanjilEndKBM = ganjilEndKBM
	existing.SemGanjilEndAssessment = ganjilEndAssesment
	existing.SemGenapStartDate = genapStart
	existing.SemGenapEndKBM = genapEndKBM
	existing.SemGenapEndAssessment = genapEndAssessment
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, existing); err != nil {
		slog.Error("Failed to update academic year", "error", err)
		return nil, ErrSystemFail
	}

	return mapEntityToDTO(existing), nil
}

func (s *service) Activate(ctx context.Context, id string) (*AcademicYearResponse, error) {
	err := s.repo.ActivateYear(ctx, id)
	if err != nil {
		if errors.Is(err, ErrAcademicYearNotFound) || errors.Is(err, ErrMultipleActiveYears) {
			return nil, err
		}
		slog.Error("Failed to activate year", "error", err, "id", id)
		return nil, ErrSystemFail
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		slog.Error("Failed to get activated year", "error", err, "id", id)
		return nil, ErrSystemFail
	}

	return mapEntityToDTO(existing), nil
}

func (s *service) UpdateSemesterStatus(ctx context.Context, id string, req SemesterStatusRequest) (*AcademicYearResponse, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrAcademicYearNotFound) {
			return nil, ErrAcademicYearNotFound
		}
		slog.Error("Failed to get year for status update", "error", err, "id", id)
		return nil, ErrSystemFail
	}

	switch req.Status {
	case SemStatusNotActive, SemStatusActive, SemStatusAssessment, SemStatusReadyToClose, SemStatusLocked:
		// Valid
	default:
		return nil, ErrInvalidSemesterStatus
	}

	if req.Semester == SemesterOdd {
		existing.SemGanjilStatus = req.Status
	} else if req.Semester == SemesterEven {
		existing.SemGenapStatus = req.Status
	}
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, existing); err != nil {
		slog.Error("Failed to update semester status", "error", err)
		return nil, ErrSystemFail
	}

	return mapEntityToDTO(existing), nil
}

func (s *service) CloseSemester(ctx context.Context, id string, req CloseSemesterRequest) (*AcademicYearResponse, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrAcademicYearNotFound) {
			return nil, ErrAcademicYearNotFound
		}
		slog.Error("Failed to get year to close semester", "error", err, "id", id)
		return nil, ErrSystemFail
	}

	if req.Semester == SemesterOdd {
		if existing.SemGanjilStatus == SemStatusNotActive || existing.SemGanjilStatus == SemStatusLocked {
			return nil, ErrSemesterNotActive
		}
		existing.SemGanjilStatus = SemStatusLocked
		existing.SemGenapStatus = SemStatusActive
	} else if req.Semester == SemesterEven {
		if existing.SemGenapStatus == SemStatusNotActive || existing.SemGenapStatus == SemStatusLocked {
			return nil, ErrSemesterNotActive
		}
		existing.SemGenapStatus = SemStatusLocked
		existing.Status = StatusPendingPromotion // Menunggu Kenaikan
	}
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, existing); err != nil {
		slog.Error("Failed to close semester", "error", err)
		return nil, ErrSystemFail
	}

	return mapEntityToDTO(existing), nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrAcademicYearNotFound) {
			return ErrAcademicYearNotFound
		}
		slog.Error("Failed to get year for delete", "error", err, "id", id)
		return ErrSystemFail
	}

	if existing.Status == StatusActive {
		return ErrDeleteActiveYear
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		slog.Error("Failed to delete academic year", "error", err)
		return ErrSystemFail
	}

	return nil
}

func mapEntityToDTO(y *AcademicYear) *AcademicYearResponse {
	activeSemesterName := "-"
	if y.Status == StatusActive {
		if y.SemGanjilStatus == SemStatusActive || y.SemGanjilStatus == SemStatusAssessment || y.SemGanjilStatus == SemStatusReadyToClose {
			activeSemesterName = "GANJIL"
		} else if y.SemGenapStatus == SemStatusActive || y.SemGenapStatus == SemStatusAssessment || y.SemGenapStatus == SemStatusReadyToClose {
			activeSemesterName = "GENAP"
		}
	}

	return &AcademicYearResponse{
		ID:        y.ID,
		Year:      y.Year,
		StartDate: formatDate(y.StartDate),
		EndDate:   formatDate(y.EndDate),
		Status:    y.Status,
		IsActive:  y.Status == StatusActive,
		Semester:  activeSemesterName,
		SemesterGanjil: SemesterData{
			StartDate:     formatNullDate(y.SemGanjilStartDate),
			EndKBM:        formatNullDate(y.SemGanjilEndKBM),
			EndAssessment: formatNullDate(y.SemGanjilEndAssessment),
			Status:        y.SemGanjilStatus,
		},
		SemesterGenap: SemesterData{
			StartDate:     formatNullDate(y.SemGenapStartDate),
			EndKBM:        formatNullDate(y.SemGenapEndKBM),
			EndAssessment: formatNullDate(y.SemGenapEndAssessment),
			Status:        y.SemGenapStatus,
		},
		CreatedAt: y.CreatedAt,
		UpdatedAt: y.UpdatedAt,
	}
}
