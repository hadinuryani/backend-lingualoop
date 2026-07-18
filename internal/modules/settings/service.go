package settings

import (
	"context"
	"log/slog"
	"time"

	"backend-lingualoop/pkg/storage"
)

type Service interface {
	GetConfig(ctx context.Context) (*SchoolSettingsResponse, error)
	UpdateConfig(ctx context.Context, req SchoolSettingsRequest) (*SchoolSettingsResponse, error)
}

type service struct {
	repo  Repository
	store storage.Storage
}

func NewService(repo Repository, store storage.Storage) Service {
	return &service{repo: repo, store: store}
}

func getURL(store storage.Storage, path *string) string {
	if path == nil || *path == "" {
		return ""
	}
	return store.GetURL(*path)
}

func mapEntityToDTO(s *SchoolSettings, store storage.Storage) *SchoolSettingsResponse {
	return &SchoolSettingsResponse{
		ID:                       s.ID,
		SchoolName:               s.SchoolName,
		SchoolNPSN:               s.SchoolNPSN,
		SchoolAddress:            s.SchoolAddress,
		SchoolPhone:              s.SchoolPhone,
		SchoolEmail:              s.SchoolEmail,
		SchoolLogoFileID:         s.SchoolLogoFileID,
		SchoolLogoURL:            getURL(store, s.SchoolLogoPath),
		EducationLogoFileID:      s.EducationLogoFileID,
		EducationLogoURL:         getURL(store, s.EducationLogoPath),
		PrincipalName:            s.PrincipalName,
		PrincipalNIP:             s.PrincipalNIP,
		PrincipalSignatureFileID: s.PrincipalSignatureFileID,
		PrincipalSignatureURL:    getURL(store, s.PrincipalSignaturePath),
		MaxStudentsPerClass:      s.MaxStudentsPerClass,
		GradingSystem:            s.GradingSystem,
		PassingGrade:             s.PassingGrade,
		AppName:                  s.AppName,
		EnableStudentLogin:       s.EnableStudentLogin,
		EnableTeacherLogin:       s.EnableTeacherLogin,
		MaintenanceMode:          s.MaintenanceMode,
		CreatedAt:                s.CreatedAt,
		UpdatedAt:                s.UpdatedAt,
	}
}

func (s *service) GetConfig(ctx context.Context) (*SchoolSettingsResponse, error) {
	c, err := s.repo.GetConfig(ctx)
	if err != nil {
		if err == ErrSettingsNotFound {
			// If not found, create a dummy struct, although DB migration handles it
			return &SchoolSettingsResponse{}, nil
		}
		slog.Error("Failed to query settings", "error", err)
		return nil, ErrSystemFail
	}
	return mapEntityToDTO(c, s.store), nil
}

func (s *service) UpdateConfig(ctx context.Context, req SchoolSettingsRequest) (*SchoolSettingsResponse, error) {
	c, err := s.repo.GetConfig(ctx)
	if err != nil && err != ErrSettingsNotFound {
		slog.Error("Failed to check existing settings", "error", err)
		return nil, ErrSystemFail
	}

	if c == nil {
		c = &SchoolSettings{}
	}

	c.SchoolName = req.SchoolName
	c.SchoolNPSN = req.SchoolNPSN
	c.SchoolAddress = req.SchoolAddress
	c.SchoolPhone = req.SchoolPhone
	c.SchoolEmail = req.SchoolEmail
	c.SchoolLogoFileID = req.SchoolLogoFileID
	c.EducationLogoFileID = req.EducationLogoFileID
	c.PrincipalName = req.PrincipalName
	c.PrincipalNIP = req.PrincipalNIP
	c.PrincipalSignatureFileID = req.PrincipalSignatureFileID
	c.MaxStudentsPerClass = req.MaxStudentsPerClass
	c.GradingSystem = req.GradingSystem
	c.PassingGrade = req.PassingGrade
	c.AppName = req.AppName
	c.EnableStudentLogin = req.EnableStudentLogin
	c.EnableTeacherLogin = req.EnableTeacherLogin
	c.MaintenanceMode = req.MaintenanceMode

	if err := s.repo.UpdateConfig(ctx, c); err != nil {
		slog.Error("Failed to save settings", "error", err)
		return nil, ErrSystemFail
	}

	// Fetch back to get updated paths and timestamps
	updated, err := s.repo.GetConfig(ctx)
	if err != nil {
		// Fallback to manual assignment if refetch fails
		c.UpdatedAt = time.Now()
		return mapEntityToDTO(c, s.store), nil
	}

	return mapEntityToDTO(updated, s.store), nil
}
