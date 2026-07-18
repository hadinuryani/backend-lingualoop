package schedule

import "errors"

var (
	ErrScheduleNotFound = errors.New("data jadwal tidak ditemukan")
	ErrClassClash       = errors.New("slot jadwal ini sudah terisi oleh mata pelajaran lain di kelas yang sama")
	ErrTeacherClash     = errors.New("guru ini sudah mengajar di kelas lain pada hari dan jam yang sama")
	ErrSystemFail       = errors.New("terjadi kesalahan pada sistem")
)
