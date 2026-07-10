package class

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
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

// Logic untuk mencari next letter
func generateClassNames(existingNames []string, levelName, majorCode string, count int) []string {
	const LETTERS = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Cari abjad tertinggi yang sudah terpakai
	nextLetterIndex := 0
	if len(existingNames) > 0 {
		var indices []int
		for _, name := range existingNames {
			parts := strings.Split(strings.TrimSpace(name), " ")
			if len(parts) > 0 {
				lastChar := parts[len(parts)-1]
				idx := strings.Index(LETTERS, strings.ToUpper(lastChar))
				if idx >= 0 {
					indices = append(indices, idx)
				}
			}
		}

		if len(indices) > 0 {
			maxIdx := -1
			for _, idx := range indices {
				if idx > maxIdx {
					maxIdx = idx
				}
			}
			nextLetterIndex = maxIdx + 1
		}
	}

	var newNames []string
	baseName := majorCode

	for i := 0; i < count; i++ {
		// Kalau abis Z, dia akan loop balik ke A (tapi sangat jarang terjadi di dunia nyata)
		letter := string(LETTERS[(nextLetterIndex+i)%26])
		className := fmt.Sprintf("%s %s %s", levelName, baseName, letter)
		newNames = append(newNames, className)
	}

	return newNames
}

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

	existingID, err := s.repo.GetIDByNameAndYear(ctx, req.ClassName, req.AcademicYearID)
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
		LevelID:           req.LevelID,
		ClassName:         req.ClassName,
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

	for _, levelID := range req.LevelIDs {
		levelName := getLevelName(levelID)

		// Tarik nama kelas yang sudah ada di DB untuk (level + major + year)
		existingNames, err := s.repo.FindNamesByLevelMajorYear(ctx, levelID, req.MajorID, req.AcademicYearID)
		if err != nil {
			slog.Error("Failed to find existing class names for batch", "error", err)
			return nil, ErrSystemFail
		}

		// Generate
		newNames := generateClassNames(existingNames, levelName, req.MajorCode, req.Count)

		for _, className := range newNames {
			c := &Class{
				ID:                uuid.New().String(),
				AcademicYearID:    req.AcademicYearID,
				MajorID:           &req.MajorID,
				LevelID:           levelID,
				ClassName:         className,
				Classroom:         nil, // Default kosong di awal
				Capacity:          capacity,
				HomeroomTeacherID: nil, // Belum ada wali kelas
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
