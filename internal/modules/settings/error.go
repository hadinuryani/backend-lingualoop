package settings

import "errors"

var (
	ErrSettingsNotFound = errors.New("pengaturan tidak ditemukan")
	ErrSystemFail       = errors.New("terjadi kesalahan pada sistem pengaturan")
)
