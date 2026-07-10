package teacher

import "errors"

var (
	ErrTeacherNotFound  = errors.New("data guru tidak ditemukan")
	ErrNipExists        = errors.New("guru dengan NIP tersebut sudah terdaftar")
	ErrInvalidStatus    = errors.New("status guru tidak valid")
	ErrInvalidBirthDate = errors.New("format tanggal lahir tidak valid")
	ErrSystemFail       = errors.New("terjadi kesalahan pada sistem, silakan coba lagi")
)
