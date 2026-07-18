package class

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	GetAll(ctx context.Context) ([]*ClassResponse, error)
	GetByID(ctx context.Context, id string) (*ClassResponse, error)
	Create(ctx context.Context, req ClassRequest) (*ClassResponse, error)
	CreateBatch(ctx context.Context, req ClassBatchRequest) ([]*ClassResponse, error)
	Update(ctx context.Context, id string, req ClassUpdateRequest) (*ClassResponse, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func optionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

var levelNames = map[string]string{
	"lvl-10": "X",
	"lvl-11": "XI",
	"lvl-12": "XII",
}

func getLevelName(levelID string) string {
	if name, ok := levelNames[levelID]; ok {
		return name
	}
	return levelID // Fallback
}

// generateClassNames logic is removed as we now generate by GradeLevel and ClassNumber

func (s *service) GetAll(ctx context.Context) ([]*ClassResponse, error) {
	classes, err := s.repo.FindAll(ctx)
	if err != nil {
		slog.Error("Failed to query classes", "error", err)
		return nil, ErrSystemFail
	}

	var responses []*ClassResponse
	for _, c := range classes {
		responses = append(responses, mapEntityToDTO(c))
	}
	if responses == nil {
		responses = []*ClassResponse{}
	}
	return responses, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*ClassResponse, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrClassNotFound) {
			return nil, ErrClassNotFound
		}
		slog.Error("Failed to query class by id", "error", err, "id", id)
		return nil, ErrSystemFail
	}
	return mapEntityToDTO(c), nil
}

func (s *service) Create(ctx context.Context, req ClassRequest) (*ClassResponse, error) {
	if req.MajorID == "" {
		return nil, ErrMajorRequired
	}

	levelID, levelName, err := s.repo.GetLevelByGrade(ctx, req.GradeLevel)
	if err != nil {
		slog.Error("Failed to get level by grade", "error", err)
		return nil, ErrSystemFail
	}

	majorCode, err := s.repo.GetMajorCodeByID(ctx, req.MajorID)
	if err != nil {
		slog.Error("Failed to get major code", "error", err)
		return nil, ErrSystemFail
	}

	className := fmt.Sprintf("%s %s %d", levelName, majorCode, req.ClassNumber)

	existingID, err := s.repo.GetIDByNameAndYear(ctx, className, req.AcademicYearID)
	if err != nil {
		slog.Error("Failed to check unique class name", "error", err)
		return nil, ErrSystemFail
	}
	if existingID != "" {
		return nil, ErrClassNameExists
	}

	now := time.Now()
	capacity := req.Capacity
	if capacity <= 0 {
		capacity = 36 // Default
	}

	newClass := &Class{
		ID:                uuid.New().String(),
		AcademicYearID:    req.AcademicYearID,
		MajorID:           optionalString(req.MajorID),
		LevelID:           levelID,
		GradeLevel:        req.GradeLevel,
		ClassNumber:       req.ClassNumber,
		ClassName:         className,
		Classroom:         optionalString(req.Classroom),
		Capacity:          capacity,
		HomeroomTeacherID: optionalString(req.HomeroomTeacherID),
		IsActive:          true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := s.repo.Create(ctx, newClass); err != nil {
		slog.Error("Failed to create class", "error", err)
		return nil, ErrSystemFail
	}

	return mapEntityToDTO(newClass), nil
}

func (s *service) CreateBatch(ctx context.Context, req ClassBatchRequest) ([]*ClassResponse, error) {
	var classesToInsert []*Class
	now := time.Now()

	capacity := req.Capacity
	if capacity <= 0 {
		capacity = 36
	}

	for _, gradeLevel := range req.GradeLevels {
		levelID, levelName, err := s.repo.GetLevelByGrade(ctx, gradeLevel)
		if err != nil {
			slog.Error("Failed to get level by grade", "error", err)
			return nil, ErrSystemFail
		}

		for i := 1; i <= req.Count; i++ {
			className := fmt.Sprintf("%s %s %d", levelName, req.MajorCode, i)
			
			c := &Class{
				ID:                uuid.New().String(),
				AcademicYearID:    req.AcademicYearID,
				MajorID:           &req.MajorID,
				LevelID:           levelID,
				GradeLevel:        gradeLevel,
				ClassNumber:       i,
				ClassName:         className,
				Classroom:         nil,
				Capacity:          capacity,
				HomeroomTeacherID: nil,
				IsActive:          true,
				CreatedAt:         now,
				UpdatedAt:         now,
			}
			classesToInsert = append(classesToInsert, c)
		}
	}

	// Insert all in one transaction
	if len(classesToInsert) > 0 {
		if err := s.repo.CreateBatch(ctx, classesToInsert); err != nil {
			slog.Error("Failed to create batch classes", "error", err)
			return nil, ErrSystemFail
		}
	}

	var responses []*ClassResponse
	for _, c := range classesToInsert {
		responses = append(responses, mapEntityToDTO(c))
	}
	if responses == nil {
		responses = []*ClassResponse{}
	}

	return responses, nil
}

func (s *service) Update(ctx context.Context, id string, req ClassUpdateRequest) (*ClassResponse, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrClassNotFound) {
			return nil, ErrClassNotFound
		}
		slog.Error("Failed to get class for update", "error", err)
		return nil, ErrSystemFail
	}

	if req.Capacity > 0 {
		existing.Capacity = req.Capacity
	}
	existing.HomeroomTeacherID = optionalString(req.HomeroomTeacherID)
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, existing); err != nil {
		if errors.Is(err, ErrClassNotFound) {
			return nil, ErrClassNotFound
		}
		slog.Error("Failed to update class", "error", err)
		return nil, ErrSystemFail
	}

	return mapEntityToDTO(existing), nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrClassNotFound) {
			return ErrClassNotFound
		}
		slog.Error("Failed to get class for delete", "error", err)
		return ErrSystemFail
	}

	if err := s.repo.Delete(ctx, existing.ID); err != nil {
		if errors.Is(err, ErrClassNotFound) {
			return ErrClassNotFound
		}
		slog.Error("Failed to delete class", "error", err)
		return ErrSystemFail
	}
	return nil
}

func mapEntityToDTO(c *Class) *ClassResponse {
	resp := &ClassResponse{
		ID:             c.ID,
		AcademicYearID: c.AcademicYearID,
		LevelID:        c.LevelID,
		ClassName:      c.ClassName,
		Capacity:       c.Capacity,
		IsActive:       c.IsActive,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
	if c.MajorID != nil {
		resp.MajorID = *c.MajorID
	}
	if c.Classroom != nil {
		resp.Classroom = *c.Classroom
	}
	if c.HomeroomTeacherID != nil {
		resp.HomeroomTeacherID = *c.HomeroomTeacherID
	}
	return resp
}
