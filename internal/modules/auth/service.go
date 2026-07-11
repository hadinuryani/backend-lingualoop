package auth

import (
	"context"
	"log/slog"

	"backend-lingualoop/pkg/jwt"
	"backend-lingualoop/pkg/security"
)

type Service interface {
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
}

type service struct {
	repo       Repository
	jwtManager jwt.Manager
}

func NewService(repo Repository, jwtManager jwt.Manager) Service {
	return &service{repo, jwtManager}
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {

	user, err := s.repo.FindByIdentifier(ctx, req.Username)
	if err != nil {
		slog.Error("Database query failed during login", "error", err, "identifier", req.Username)
		return nil, ErrSystemFail
	}

	if user == nil {
		return nil, ErrInvalidCreds
	}

	if !user.IsActive {
		return nil, ErrAccountDisabled
	}

	if !security.CheckPassword(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCreds
	}

	token, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username, user.Role)
	if err != nil {
		return nil, ErrSessionFail
	}

	userDTO := UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FullName:  user.FullName,
		Role:      user.Role,
		AvatarURL: user.AvatarURL,
		IsActive:  user.IsActive,
	}

	response := &LoginResponse{
		User:  userDTO,
		Token: token,
	}

	return response, nil
}
