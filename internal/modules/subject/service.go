package subject

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	GetAll(ctx context.Context) ([]*SubjectResponse, error)
	GetByID(ctx context.Context, id string) (*SubjectResponse, error)
	Create(ctx context.Context, req SubjectRequest) (*SubjectResponse, error)
	Update(ctx context.Context, id string, req SubjectRequest) (*SubjectResponse, error)
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

func (s *service) GetAll(ctx context.Context) ([]*SubjectResponse, error) {
	subjects, err := s.repo.FindAll(ctx)
	if err != nil {
		slog.Error("Failed to query subjects", "error", err)
		return nil, ErrSystemFail
	}

	var responses []*SubjectResponse
	for _, sub := range subjects {
		responses = append(responses, mapEntityToDTO(sub))
	}

	if responses == nil {
		responses = []*SubjectResponse{}
	}

	return responses, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*SubjectResponse, error) {
	sub, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrSubjectNotFound) {
			return nil, ErrSubjectNotFound
		}
		slog.Error("Failed to query subject by ID", "error", err, "id", id)
		return nil, ErrSystemFail
	}
	return mapEntityToDTO(sub), nil
}

func (s *service) Create(ctx context.Context, req SubjectRequest) (*SubjectResponse, error) {
	req.Code = strings.ToUpper(strings.TrimSpace(req.Code))

	existingID, err := s.repo.GetIDByCode(ctx, req.Code)
	if err != nil {
		slog.Error("Failed to check subject code unique", "error", err, "code", req.Code)
		return nil, ErrSystemFail
	}
	if existingID != "" {
		return nil, ErrSubjectCodeExists
	}

	now := time.Now()
	newSub := &Subject{
		ID:          uuid.New().String(),
		Code:        req.Code,
		Name:        req.Name,
		Description: optionalString(req.Description),
		MajorID:     optionalString(req.MajorID),
		LevelID:     optionalString(req.LevelID),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, newSub); err != nil {
		slog.Error("Failed to create subject", "error", err)
		return nil, ErrSystemFail
	}

	return mapEntityToDTO(newSub), nil
}

func (s *service) Update(ctx context.Context, id string, req SubjectRequest) (*SubjectResponse, error) {
	req.Code = strings.ToUpper(strings.TrimSpace(req.Code))

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrSubjectNotFound) {
			return nil, ErrSubjectNotFound
		}
		slog.Error("Failed to find subject for update", "error", err, "id", id)
		return nil, ErrSystemFail
	}

	if req.Code != existing.Code {
		existingID, err := s.repo.GetIDByCode(ctx, req.Code)
		if err != nil {
			slog.Error("Failed to check subject code unique on update", "error", err, "code", req.Code)
			return nil, ErrSystemFail
		}
		if existingID != "" && existingID != id {
			return nil, ErrSubjectCodeExists
		}
	}

	existing.Code = req.Code
	existing.Name = req.Name
	existing.Description = optionalString(req.Description)
	existing.MajorID = optionalString(req.MajorID)
	existing.LevelID = optionalString(req.LevelID)
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, existing); err != nil {
		if errors.Is(err, ErrSubjectNotFound) {
			return nil, ErrSubjectNotFound
		}
		slog.Error("Failed to update subject", "error", err, "id", id)
		return nil, ErrSystemFail
	}

	return mapEntityToDTO(existing), nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrSubjectNotFound) {
			return ErrSubjectNotFound
		}
		slog.Error("Failed to find subject for delete", "error", err, "id", id)
		return ErrSystemFail
	}

	if err := s.repo.Delete(ctx, existing.ID); err != nil {
		if errors.Is(err, ErrSubjectNotFound) {
			return ErrSubjectNotFound
		}
		slog.Error("Failed to delete subject", "error", err, "id", id)
		return ErrSystemFail
	}

	return nil
}

func mapEntityToDTO(sub *Subject) *SubjectResponse {
	resp := &SubjectResponse{
		ID:        sub.ID,
		Code:      sub.Code,
		Name:      sub.Name,
		CreatedAt: sub.CreatedAt,
		UpdatedAt: sub.UpdatedAt,
	}
	if sub.Description != nil {
		resp.Description = *sub.Description
	}
	if sub.MajorID != nil {
		resp.MajorID = *sub.MajorID
	}
	if sub.LevelID != nil {
		resp.LevelID = *sub.LevelID
	}
	return resp
}
