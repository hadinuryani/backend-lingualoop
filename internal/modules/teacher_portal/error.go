package teacher_portal

import "errors"

var (
	ErrTeacherNotFound = errors.New("guru tidak ditemukan untuk user ini")
	ErrNoActiveAcademicYear = errors.New("tidak ada tahun ajaran aktif")
)
