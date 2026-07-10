package teacher

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"backend-lingualoop/pkg/security"
	"github.com/google/uuid"
)

type Service interface {
	GetAll(ctx context.Context) ([]*TeacherResponse, error)
	GetByID(ctx context.Context, id string) (*TeacherResponse, error)
	Create(ctx context.Context, req TeacherRequest) (*TeacherResponse, error)
	Update(ctx context.Context, id string, req TeacherRequest) (*TeacherResponse, error)
	ToggleStatus(ctx context.Context, id string) (*TeacherResponse, error)
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

// Konfigurasi sementara sebelum tabel settings ada
const DefaultTeacherPassword = "guru123"

func optionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func generateTeacherUsername(fullName, nip string) string {
	firstName := strings.ToLower(strings.Split(fullName, " ")[0])
	var cleanFirstName strings.Builder
	for _, r := range firstName {
		if r >= 'a' && r <= 'z' {
			cleanFirstName.WriteRune(r)
		}
	}
	if cleanFirstName.Len() == 0 {
		cleanFirstName.WriteString("user")
	}

	nipSuffix := nip
	if len(nip) > 3 {
		nipSuffix = nip[len(nip)-3:]
	}

	return fmt.Sprintf("%s%s", cleanFirstName.String(), nipSuffix)
}

func generateTeacherEmail(username string) string {
	return fmt.Sprintf("%s@teacher.lingualoop.id", username)
}

func buildUser(email, username, passwordHash, fullName string, now time.Time) *User {
	return &User{
		ID:           uuid.New().String(),
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
		FullName:     fullName,
		Role:         RoleTeacher,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func buildTeacher(userID, username, email string, req TeacherRequest, now time.Time) (*Teacher, error) {
	t := &Teacher{
		ID:            uuid.New().String(),
		UserID:        userID,
		Username:      username,
		Email:         email,
		NIP:           req.NIP,
		FullName:      req.FullName,
		Gender:        req.Gender,
		Status:        TeacherActive,
		BirthPlace:    optionalString(req.BirthPlace),
		Phone:         optionalString(req.Phone),
		AddressRegion: optionalString(req.AddressRegion),
		AddressDetail: optionalString(req.AddressDetail),
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if req.BirthDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			return nil, ErrInvalidBirthDate
		}
		t.BirthDate = &parsedDate
	}

	return t, nil
}

func (s *service) GetAll(ctx context.Context) ([]*TeacherResponse, error) {
	teachers, err := s.repo.FindAll(ctx)
	if err != nil {
		slog.Error("Failed to query teachers", "error", err)
		return nil, ErrSystemFail
	}

	var responses []*TeacherResponse
	for _, t := range teachers {
		responses = append(responses, mapEntityToDTO(t))
	}

	if responses == nil {
		responses = []*TeacherResponse{}
	}

	return responses, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*TeacherResponse, error) {
	teacher, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrTeacherNotFound) {
			return nil, ErrTeacherNotFound
		}
		slog.Error("Failed to query teacher by ID", "error", err, "id", id)
		return nil, ErrSystemFail
	}
	return mapEntityToDTO(teacher), nil
}

func (s *service) Create(ctx context.Context, req TeacherRequest) (*TeacherResponse, error) {
	existingID, err := s.repo.GetIDByNIP(ctx, req.NIP)
	if err != nil {
		slog.Error("Failed to query teacher ID by NIP", "error", err, "nip", req.NIP)
		return nil, ErrSystemFail
	}
	if existingID != "" {
		return nil, ErrNipExists
	}

	username := generateTeacherUsername(req.FullName, req.NIP)
	email := generateTeacherEmail(username)

	passwordHash, err := security.HashPassword(DefaultTeacherPassword)
	if err != nil {
		slog.Error("Failed to hash default password", "error", err)
		return nil, ErrSystemFail
	}

	now := time.Now()

	newUser := buildUser(email, username, passwordHash, req.FullName, now)

	newTeacher, err := buildTeacher(newUser.ID, username, email, req, now)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, newTeacher, newUser); err != nil {
		slog.Error("Failed to create teacher and user", "error", err, "nip", req.NIP)
		return nil, ErrSystemFail
	}

	return mapEntityToDTO(newTeacher), nil
}

func (s *service) Update(ctx context.Context, id string, req TeacherRequest) (*TeacherResponse, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrTeacherNotFound) {
			return nil, ErrTeacherNotFound
		}
		slog.Error("Failed to find teacher for update", "error", err, "id", id)
		return nil, ErrSystemFail
	}

	if req.NIP != existing.NIP {
		existingID, err := s.repo.GetIDByNIP(ctx, req.NIP)
		if err != nil {
			slog.Error("Failed to check nip unique on update", "error", err, "nip", req.NIP)
			return nil, ErrSystemFail
		}
		if existingID != "" && existingID != id {
			return nil, ErrNipExists
		}
	}

	existing.NIP = req.NIP
	existing.FullName = req.FullName
	existing.Gender = req.Gender
	existing.UpdatedAt = time.Now()

	existing.BirthPlace = optionalString(req.BirthPlace)
	existing.Phone = optionalString(req.Phone)
	existing.AddressRegion = optionalString(req.AddressRegion)
	existing.AddressDetail = optionalString(req.AddressDetail)

	if req.BirthDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			return nil, ErrInvalidBirthDate
		}
		existing.BirthDate = &parsedDate
	} else {
		existing.BirthDate = nil
	}

	userToUpdate := &User{
		ID:        existing.UserID,
		FullName:  req.FullName,
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Update(ctx, existing, userToUpdate); err != nil {
		if errors.Is(err, ErrTeacherNotFound) {
			return nil, ErrTeacherNotFound
		}
		slog.Error("Failed to update teacher", "error", err, "id", id)
		return nil, ErrSystemFail
	}

	return mapEntityToDTO(existing), nil
}

func (s *service) ToggleStatus(ctx context.Context, id string) (*TeacherResponse, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrTeacherNotFound) {
			return nil, ErrTeacherNotFound
		}
		slog.Error("Failed to find teacher for status toggle", "error", err, "id", id)
		return nil, ErrSystemFail
	}

	newStatus := TeacherActive
	if existing.Status == TeacherActive {
		newStatus = TeacherInactive
	}

	if err := s.repo.UpdateStatus(ctx, existing.ID, existing.UserID, newStatus); err != nil {
		if errors.Is(err, ErrTeacherNotFound) {
			return nil, ErrTeacherNotFound
		}
		if errors.Is(err, ErrInvalidStatus) {
			return nil, ErrInvalidStatus
		}
		slog.Error("Failed to toggle teacher status", "error", err, "id", id)
		return nil, ErrSystemFail
	}

	existing.Status = newStatus
	return mapEntityToDTO(existing), nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrTeacherNotFound) {
			return ErrTeacherNotFound
		}
		slog.Error("Failed to find teacher for delete", "error", err, "id", id)
		return ErrSystemFail
	}

	if err := s.repo.Delete(ctx, existing.ID, existing.UserID); err != nil {
		if errors.Is(err, ErrTeacherNotFound) {
			return ErrTeacherNotFound
		}
		slog.Error("Failed to delete teacher", "error", err, "id", id)
		return ErrSystemFail
	}

	return nil
}

func mapEntityToDTO(t *Teacher) *TeacherResponse {
	resp := &TeacherResponse{
		ID:        t.ID,
		UserID:    t.UserID,
		Username:  t.Username,
		NIP:       t.NIP,
		FullName:  t.FullName,
		Gender:    t.Gender,
		Status:    t.Status,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
	if t.BirthPlace != nil {
		resp.BirthPlace = *t.BirthPlace
	}
	if t.BirthDate != nil {
		resp.BirthDate = t.BirthDate.Format("2006-01-02")
	}
	if t.Phone != nil {
		resp.Phone = *t.Phone
	}
	if t.AddressRegion != nil {
		resp.AddressRegion = *t.AddressRegion
	}
	if t.AddressDetail != nil {
		resp.AddressDetail = *t.AddressDetail
	}
	if t.Photo != nil {
		resp.Photo = *t.Photo
	}
	return resp
}
