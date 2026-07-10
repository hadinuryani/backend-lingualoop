package subject

import "errors"

var (
	ErrSubjectNotFound   = errors.New("data mata pelajaran tidak ditemukan")
	ErrSubjectCodeExists = errors.New("kode mata pelajaran sudah terdaftar")
	ErrSystemFail        = errors.New("terjadi kesalahan pada sistem, silakan coba lagi")
)
