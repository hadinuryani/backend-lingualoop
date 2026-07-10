package class

import "errors"

var (
	ErrClassNotFound   = errors.New("data kelas tidak ditemukan")
	ErrClassNameExists = errors.New("nama kelas sudah terdaftar pada tahun akademik ini")
	ErrSystemFail      = errors.New("terjadi kesalahan pada sistem, silakan coba lagi")
	ErrInvalidLevel    = errors.New("tingkatan kelas tidak valid")
	ErrMajorRequired   = errors.New("jurusan wajib diisi")
)
