package auth

import "errors"

var (
	ErrSystemFail      = errors.New("terjadi kesalahan pada sistem, silakan coba lagi")
	ErrInvalidCreds    = errors.New("email atau password salah")
	ErrAccountDisabled = errors.New("akun Anda dinonaktifkan, silakan hubungi admin")
	ErrSessionFail     = errors.New("gagal membuat sesi login")
)
