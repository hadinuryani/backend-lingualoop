package student

import "errors"

var (
	ErrStudentNotFound  = errors.New("data siswa tidak ditemukan")
	ErrNisExists        = errors.New("siswa dengan NIS tersebut sudah terdaftar")
	ErrInvalidStatus    = errors.New("status siswa tidak valid")
	ErrInvalidBirthDate = errors.New("format tanggal lahir tidak valid")
	ErrSystemFail       = errors.New("terjadi kesalahan pada sistem, silakan coba lagi")
)
