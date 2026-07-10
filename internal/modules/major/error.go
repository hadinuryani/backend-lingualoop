package major

import "errors"

var (
	ErrSystemFail    = errors.New("terjadi kesalahan pada sistem, silakan coba lagi")
	ErrCodeExists    = errors.New("jurusan dengan kode tersebut sudah terdaftar")
	ErrNameExists    = errors.New("jurusan dengan nama tersebut sudah terdaftar")
	ErrMajorNotFound = errors.New("data jurusan tidak ditemukan")
)
