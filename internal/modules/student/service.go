package student

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
	GetAll(ctx context.Context) ([]*StudentResponse, error)
	GetByID(ctx context.Context, id string) (*StudentResponse, error)
	Create(ctx context.Context, req StudentRequest) (*StudentResponse, error)
	Update(ctx context.Context, id string, req StudentRequest) (*StudentResponse, error)
	UpdateStatus(ctx context.Context, id string, status string) (*StudentResponse, error)
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

const DefaultStudentPassword = "siswa123"

func optionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func generateStudentUsername(fullName, nis string) string {
	firstName := strings.ToLower(strings.Split(fullName, " ")[0])
	var cleanFirstName strings.Builder
	for _, r := range firstName {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			cleanFirstName.WriteRune(r)
		}
	}
	if cleanFirstName.Len() == 0 {
		cleanFirstName.WriteString("user")
	}

	nisSuffix := nis
	if len(nis) > 3 {
		nisSuffix = nis[len(nis)-3:]
	}

	return fmt.Sprintf("%s%s", cleanFirstName.String(), nisSuffix)
}

func generateStudentEmail(username string) string {
	return fmt.Sprintf("%s@student.lingualoop.id", username)
}

func buildUser(email, username, passwordHash, fullName string, now time.Time) *User {
	return &User{
		ID:           uuid.New().String(),
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
		FullName:     fullName,
		Role:         RoleStudent,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func buildStudent(userID, username, email string, req StudentRequest, now time.Time) (*Student, error) {
	s := &Student{
		ID:            uuid.New().String(),
		UserID:        userID,
		Username:      username,
		Email:         email,
		NIS:           req.NIS,
		FullName:      req.FullName,
		Gender:        req.Gender,
		Status:        StudentActive,
		MajorID:       optionalString(req.MajorID),
		ClassLevel:    optionalString(req.ClassLevel),
		BirthPlace:    optionalString(req.BirthPlace),
		Phone:         optionalString(req.Phone),
		AddressRegion: optionalString(req.AddressRegion),
		AddressDetail: optionalString(req.AddressDetail),
		Photo:         optionalString(req.Photo),
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if req.BirthDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			return nil, ErrInvalidBirthDate
		}
		s.BirthDate = &parsedDate
	}

	return s, nil
}

func (s *service) GetAll(ctx context.Context) ([]*StudentResponse, error) {
	students, err := s.repo.FindAll(ctx)
	if err != nil {
		slog.Error("Failed to query students", "error", err)
		return nil, ErrSystemFail
	}

	var responses []*StudentResponse
	for _, std := range students {
		responses = append(responses, mapEntityToDTO(std))
	}

	if responses == nil {
		responses = []*StudentResponse{}
	}

	return responses, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*StudentResponse, error) {
	student, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			return nil, ErrStudentNotFound
		}
		slog.Error("Failed to query student by ID", "error", err, "id", id)
		return nil, ErrSystemFail
	}
	return mapEntityToDTO(student), nil
}

func (s *service) Create(ctx context.Context, req StudentRequest) (*StudentResponse, error) {
	existingID, err := s.repo.GetIDByNIS(ctx, req.NIS)
	if err != nil {
		slog.Error("Failed to query student ID by NIS", "error", err, "nis", req.NIS)
		return nil, ErrSystemFail
	}
	if existingID != "" {
		return nil, ErrNisExists
	}

	username := generateStudentUsername(req.FullName, req.NIS)
	email := generateStudentEmail(username)

	passwordHash, err := security.HashPassword(DefaultStudentPassword)
	if err != nil {
		slog.Error("Failed to hash default password", "error", err)
		return nil, ErrSystemFail
	}

	now := time.Now()

	newUser := buildUser(email, username, passwordHash, req.FullName, now)

	newStudent, err := buildStudent(newUser.ID, username, email, req, now)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, newStudent, newUser); err != nil {
		slog.Error("Failed to create student and user", "error", err, "nis", req.NIS)
		return nil, ErrSystemFail
	}

	return mapEntityToDTO(newStudent), nil
}

func (s *service) Update(ctx context.Context, id string, req StudentRequest) (*StudentResponse, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			return nil, ErrStudentNotFound
		}
		slog.Error("Failed to find student for update", "error", err, "id", id)
		return nil, ErrSystemFail
	}

	if req.NIS != existing.NIS {
		existingID, err := s.repo.GetIDByNIS(ctx, req.NIS)
		if err != nil {
			slog.Error("Failed to check nis unique on update", "error", err, "nis", req.NIS)
			return nil, ErrSystemFail
		}
		if existingID != "" && existingID != id {
			return nil, ErrNisExists
		}
	}

	existing.NIS = req.NIS
	existing.FullName = req.FullName
	existing.Gender = req.Gender
	existing.MajorID = optionalString(req.MajorID)
	existing.ClassLevel = optionalString(req.ClassLevel)
	existing.BirthPlace = optionalString(req.BirthPlace)
	existing.Phone = optionalString(req.Phone)
	existing.AddressRegion = optionalString(req.AddressRegion)
	existing.AddressDetail = optionalString(req.AddressDetail)
	existing.Photo = optionalString(req.Photo)
	existing.UpdatedAt = time.Now()

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
		if errors.Is(err, ErrStudentNotFound) {
			return nil, ErrStudentNotFound
		}
		slog.Error("Failed to update student", "error", err, "id", id)
		return nil, ErrSystemFail
	}

	return mapEntityToDTO(existing), nil
}

func (s *service) UpdateStatus(ctx context.Context, id string, status string) (*StudentResponse, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			return nil, ErrStudentNotFound
		}
		slog.Error("Failed to find student for status update", "error", err, "id", id)
		return nil, ErrSystemFail
	}

	if err := s.repo.UpdateStatus(ctx, existing.ID, existing.UserID, status); err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			return nil, ErrStudentNotFound
		}
		if errors.Is(err, ErrInvalidStatus) {
			return nil, ErrInvalidStatus
		}
		slog.Error("Failed to update student status", "error", err, "id", id)
		return nil, ErrSystemFail
	}

	existing.Status = status
	return mapEntityToDTO(existing), nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			return ErrStudentNotFound
		}
		slog.Error("Failed to find student for delete", "error", err, "id", id)
		return ErrSystemFail
	}

	if err := s.repo.Delete(ctx, existing.ID, existing.UserID); err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			return ErrStudentNotFound
		}
		slog.Error("Failed to delete student", "error", err, "id", id)
		return ErrSystemFail
	}

	return nil
}

func mapEntityToDTO(s *Student) *StudentResponse {
	resp := &StudentResponse{
		ID:        s.ID,
		UserID:    s.UserID,
		Username:  s.Username,
		NIS:       s.NIS,
		FullName:  s.FullName,
		Gender:    s.Gender,
		Status:    s.Status,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
	if s.BirthPlace != nil {
		resp.BirthPlace = *s.BirthPlace
	}
	if s.BirthDate != nil {
		resp.BirthDate = s.BirthDate.Format("2006-01-02")
	}
	if s.Phone != nil {
		resp.Phone = *s.Phone
	}
	if s.AddressRegion != nil {
		resp.AddressRegion = *s.AddressRegion
	}
	if s.AddressDetail != nil {
		resp.AddressDetail = *s.AddressDetail
	}
	if s.Photo != nil {
		resp.Photo = *s.Photo
	}
	if s.MajorID != nil {
		resp.MajorID = *s.MajorID
	}
	if s.ClassLevel != nil {
		resp.ClassLevel = *s.ClassLevel
	}
	return resp
}
