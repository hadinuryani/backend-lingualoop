package auth

import (
	"errors"

	"backend-lingualoop/pkg/jwt"
	"backend-lingualoop/pkg/security"
)

type Service interface {
	Login(req LoginRequest) (*LoginResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Login(req LoginRequest) (*LoginResponse, error) {

	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("terjadi kesalahan pada sistem, silakan coba lagi")
	}

	if user == nil {
		return nil, errors.New("email atau password salah")
	}

	if !user.IsActive {
		return nil, errors.New("akun Anda dinonaktifkan, silakan hubungi admin")
	}

	if !security.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("email atau password salah")
	}

	token, err := jwt.GenerateToken(user.ID, user.Email, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("gagal membuat sesi login")
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
