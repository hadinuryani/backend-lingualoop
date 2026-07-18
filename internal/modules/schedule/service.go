package schedule

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	GetConfig(ctx context.Context) (*ScheduleConfigResponse, error)
	SaveConfig(ctx context.Context, req ScheduleConfigRequest) (*ScheduleConfigResponse, error)

	GetAll(ctx context.Context) ([]*ScheduleResponse, error)
	GetByClass(ctx context.Context, classID, academicYearID string) ([]*ScheduleResponse, error)
	Create(ctx context.Context, req ScheduleRequest) (*ScheduleResponse, error)
	Update(ctx context.Context, id string, req ScheduleRequest) (*ScheduleResponse, error)
	Delete(ctx context.Context, id string) error
	DeleteByClass(ctx context.Context, classID, academicYearID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func mapConfigEntityToDTO(c *ScheduleConfig) *ScheduleConfigResponse {
	var breakAfterPeriods []int
	var breakDurations []int
	var activeDays []string

	if err := json.Unmarshal(c.BreakAfterPeriods, &breakAfterPeriods); err != nil {
		slog.Warn("Failed to unmarshal BreakAfterPeriods", "error", err)
	}
	if err := json.Unmarshal(c.BreakDurations, &breakDurations); err != nil {
		slog.Warn("Failed to unmarshal BreakDurations", "error", err)
	}
	if err := json.Unmarshal(c.ActiveDays, &activeDays); err != nil {
		slog.Warn("Failed to unmarshal ActiveDays", "error", err)
	}

	return &ScheduleConfigResponse{
		PeriodsPerDay:     c.PeriodsPerDay,
		PeriodDuration:    c.PeriodDuration,
		StartTime:         c.StartTime,
		BreakAfterPeriods: breakAfterPeriods,
		BreakDurations:    breakDurations,
		ActiveDays:        activeDays,
		UpdatedAt:         c.UpdatedAt,
	}
}

func mapEntityToDTO(s *Schedule) *ScheduleResponse {
	return &ScheduleResponse{
		ID:             s.ID,
		AcademicYearID: s.AcademicYearID,
		ClassID:        s.ClassID,
		SubjectID:      s.SubjectID,
		TeacherID:      s.TeacherID,
		Day:            s.Day,
		Period:         s.Period,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
	}
}

func (s *service) GetConfig(ctx context.Context) (*ScheduleConfigResponse, error) {
	c, err := s.repo.GetConfig(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &ScheduleConfigResponse{
				PeriodsPerDay:     10,
				PeriodDuration:    45,
				StartTime:         "07:00",
				BreakAfterPeriods: []int{3, 6},
				BreakDurations:    []int{15, 15},
				ActiveDays:        []string{"Senin", "Selasa", "Rabu", "Kamis", "Jumat"},
			}, nil
		}
		slog.Error("Failed to query schedule config", "error", err)
		return nil, ErrSystemFail
	}
	return mapConfigEntityToDTO(c), nil
}

func (s *service) SaveConfig(ctx context.Context, req ScheduleConfigRequest) (*ScheduleConfigResponse, error) {
	bap, err := json.Marshal(req.BreakAfterPeriods)
	if err != nil {
		return nil, ErrSystemFail
	}
	bd, err := json.Marshal(req.BreakDurations)
	if err != nil {
		return nil, ErrSystemFail
	}
	ad, err := json.Marshal(req.ActiveDays)
	if err != nil {
		return nil, ErrSystemFail
	}

	c := &ScheduleConfig{
		PeriodsPerDay:     req.PeriodsPerDay,
		PeriodDuration:    req.PeriodDuration,
		StartTime:         req.StartTime,
		BreakAfterPeriods: json.RawMessage(bap),
		BreakDurations:    json.RawMessage(bd),
		ActiveDays:        json.RawMessage(ad),
	}

	if err := s.repo.SaveConfig(ctx, c); err != nil {
		slog.Error("Failed to save schedule config", "error", err)
		return nil, ErrSystemFail
	}
	
	return &ScheduleConfigResponse{
		PeriodsPerDay:     req.PeriodsPerDay,
		PeriodDuration:    req.PeriodDuration,
		StartTime:         req.StartTime,
		BreakAfterPeriods: req.BreakAfterPeriods,
		BreakDurations:    req.BreakDurations,
		ActiveDays:        req.ActiveDays,
		UpdatedAt:         time.Now(),
	}, nil
}

func (s *service) GetAll(ctx context.Context) ([]*ScheduleResponse, error) {
	schedules, err := s.repo.FindAll(ctx)
	if err != nil {
		slog.Error("Failed to get all schedules", "error", err)
		return nil, ErrSystemFail
	}

	var res []*ScheduleResponse
	for _, sch := range schedules {
		res = append(res, mapEntityToDTO(sch))
	}
	if res == nil {
		res = []*ScheduleResponse{}
	}
	return res, nil
}

func (s *service) GetByClass(ctx context.Context, classID, academicYearID string) ([]*ScheduleResponse, error) {
	schedules, err := s.repo.FindByClass(ctx, classID, academicYearID)
	if err != nil {
		slog.Error("Failed to get schedules by class", "error", err)
		return nil, ErrSystemFail
	}

	var res []*ScheduleResponse
	for _, sch := range schedules {
		res = append(res, mapEntityToDTO(sch))
	}
	if res == nil {
		res = []*ScheduleResponse{}
	}
	return res, nil
}

func (s *service) Create(ctx context.Context, req ScheduleRequest) (*ScheduleResponse, error) {
	if req.Period <= 0 {
		return nil, errors.New("period harus lebih dari 0")
	}
	if req.Day == "" {
		return nil, errors.New("hari tidak boleh kosong")
	}

	// Validasi Class Clash
	clashClass, err := s.repo.FindClassClash(ctx, req.ClassID, req.AcademicYearID, req.Day, req.Period, "")
	if err != nil {
		slog.Error("Failed to check class clash", "error", err)
		return nil, ErrSystemFail
	}
	if clashClass {
		return nil, ErrClassClash
	}

	// Validasi Teacher Clash
	clashTeacher, err := s.repo.FindTeacherClash(ctx, req.TeacherID, req.AcademicYearID, req.Day, req.Period, "")
	if err != nil {
		slog.Error("Failed to check teacher clash", "error", err)
		return nil, ErrSystemFail
	}
	if clashTeacher {
		return nil, ErrTeacherClash
	}

	sch := &Schedule{
		ID:             uuid.New().String(),
		AcademicYearID: req.AcademicYearID,
		ClassID:        req.ClassID,
		SubjectID:      req.SubjectID,
		TeacherID:      req.TeacherID,
		Day:            req.Day,
		Period:         req.Period,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.Create(ctx, sch); err != nil {
		slog.Error("Failed to create schedule", "error", err)
		return nil, ErrSystemFail
	}
	return mapEntityToDTO(sch), nil
}

func (s *service) Update(ctx context.Context, id string, req ScheduleRequest) (*ScheduleResponse, error) {
	if req.Period <= 0 {
		return nil, errors.New("period harus lebih dari 0")
	}
	if req.Day == "" {
		return nil, errors.New("hari tidak boleh kosong")
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrScheduleNotFound) {
			return nil, ErrScheduleNotFound
		}
		slog.Error("Failed to check existing schedule", "error", err)
		return nil, ErrSystemFail
	}

	// Cek class clash jika day/period/class berubah
	if existing.Day != req.Day || existing.Period != req.Period || existing.ClassID != req.ClassID {
		clashClass, err := s.repo.FindClassClash(ctx, req.ClassID, req.AcademicYearID, req.Day, req.Period, id)
		if err != nil {
			return nil, ErrSystemFail
		}
		if clashClass {
			return nil, ErrClassClash
		}
	}

	// Cek teacher clash jika day/period/teacher berubah
	if existing.Day != req.Day || existing.Period != req.Period || existing.TeacherID != req.TeacherID {
		clashTeacher, err := s.repo.FindTeacherClash(ctx, req.TeacherID, req.AcademicYearID, req.Day, req.Period, id)
		if err != nil {
			return nil, ErrSystemFail
		}
		if clashTeacher {
			return nil, ErrTeacherClash
		}
	}

	existing.SubjectID = req.SubjectID
	existing.TeacherID = req.TeacherID
	existing.Day = req.Day
	existing.Period = req.Period
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, existing); err != nil {
		if errors.Is(err, ErrScheduleNotFound) {
			return nil, ErrScheduleNotFound
		}
		slog.Error("Failed to update schedule", "error", err)
		return nil, ErrSystemFail
	}
	
	// Refetch to get updated timestamp from DB (optional, but good)
	return mapEntityToDTO(existing), nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil && !errors.Is(err, ErrScheduleNotFound) {
		slog.Error("Failed to delete schedule", "error", err)
		return ErrSystemFail
	}
	return err
}

func (s *service) DeleteByClass(ctx context.Context, classID, academicYearID string) error {
	err := s.repo.DeleteByClass(ctx, classID, academicYearID)
	if err != nil {
		slog.Error("Failed to delete schedule by class", "error", err)
		return ErrSystemFail
	}
	return nil
}
