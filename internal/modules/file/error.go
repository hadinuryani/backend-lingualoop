package file

import "errors"

var (
	ErrFileNotFound    = errors.New("file tidak ditemukan")
	ErrInvalidFileType = errors.New("tipe file tidak valid atau ekstensi dimanipulasi")
	ErrFileTooLarge    = errors.New("ukuran file melebihi batas maksimal")
	ErrSystemFail      = errors.New("terjadi kesalahan pada sistem")
)
