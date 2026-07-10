package major

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	GetAll(ctx context.Context) ([]*MajorResponse, error)
	Create(ctx context.Context, req MajorRequest) (*MajorResponse, error)
	Update(ctx context.Context, id string, req MajorRequest) (*MajorResponse, error)
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

func (s *service) GetAll(ctx context.Context) ([]*MajorResponse, error) {
	majors, err := s.repo.FindAll(ctx)
	if err != nil {
		slog.Error("Failed to query majors", "error", err)
		return nil, ErrSystemFail
	}

	var responses []*MajorResponse
	for _, m := range majors {
		responses = append(responses, mapEntityToDTO(m))
	}

	// Pastikan return array kosong jika tidak ada data, bukan null
	if responses == nil {
		responses = []*MajorResponse{}
	}

	return responses, nil
}

func (s *service) Create(ctx context.Context, req MajorRequest) (*MajorResponse, error) {
	// Validasi kode unik
	existingCode, err := s.repo.FindByCode(ctx, strings.ToUpper(req.Code))
	if err != nil && !errors.Is(err, ErrMajorNotFound) {
		slog.Error("Failed to query major by code", "error", err, "code", req.Code)
		return nil, ErrSystemFail
	}
	if existingCode != nil {
		return nil, ErrCodeExists
	}

	// Validasi nama unik
	existingName, err := s.repo.FindByName(ctx, req.Name)
	if err != nil && !errors.Is(err, ErrMajorNotFound) {
		slog.Error("Failed to query major by name", "error", err, "name", req.Name)
		return nil, ErrSystemFail
	}
	if existingName != nil {
		return nil, ErrNameExists
	}

	now := time.Now()
	newMajor := &Major{
		ID:          uuid.New().String(),
		Code:        strings.ToUpper(req.Code),
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, newMajor); err != nil {
		slog.Error("Failed to insert major", "error", err, "code", req.Code)
		return nil, ErrSystemFail
	}

	return mapEntityToDTO(newMajor), nil
}

func (s *service) Update(ctx context.Context, id string, req MajorRequest) (*MajorResponse, error) {
	// Cek apakah data ada
	existingMajor, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrMajorNotFound) {
			return nil, ErrMajorNotFound
		}
		slog.Error("Failed to find major by ID", "error", err, "id", id)
		return nil, ErrSystemFail
	}

	// Validasi kode unik jika diubah
	if strings.ToUpper(req.Code) != existingMajor.Code {
		checkCode, err := s.repo.FindByCode(ctx, strings.ToUpper(req.Code))
		if err != nil && !errors.Is(err, ErrMajorNotFound) {
			slog.Error("Failed to check code unique on update", "error", err, "code", req.Code)
			return nil, ErrSystemFail
		}
		if checkCode != nil && checkCode.ID != id {
			return nil, ErrCodeExists
		}
	}

	// Validasi nama unik jika diubah
	if req.Name != existingMajor.Name {
		checkName, err := s.repo.FindByName(ctx, req.Name)
		if err != nil && !errors.Is(err, ErrMajorNotFound) {
			slog.Error("Failed to check name unique on update", "error", err, "name", req.Name)
			return nil, ErrSystemFail
		}
		if checkName != nil && checkName.ID != id {
			return nil, ErrNameExists
		}
	}

	existingMajor.Code = strings.ToUpper(req.Code)
	existingMajor.Name = req.Name
	existingMajor.Description = req.Description
	existingMajor.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, existingMajor); err != nil {
		if errors.Is(err, ErrMajorNotFound) {
			return nil, ErrMajorNotFound
		}
		slog.Error("Failed to update major", "error", err, "id", id)
		return nil, ErrSystemFail
	}

	return mapEntityToDTO(existingMajor), nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, ErrMajorNotFound) {
			return ErrMajorNotFound
		}
		slog.Error("Failed to delete major", "error", err, "id", id)
		return ErrSystemFail
	}

	return nil
}

// Helper mapper
func mapEntityToDTO(m *Major) *MajorResponse {
	return &MajorResponse{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
