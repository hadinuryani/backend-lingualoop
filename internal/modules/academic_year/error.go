package academic_year

import "errors"

var (
	ErrAcademicYearNotFound   = errors.New("data tahun akademik tidak ditemukan")
	ErrAcademicYearExists     = errors.New("tahun akademik sudah terdaftar")
	ErrSystemFail             = errors.New("terjadi kesalahan pada sistem, silakan coba lagi")
	ErrMultipleActiveYears    = errors.New("tidak boleh ada lebih dari satu tahun akademik yang aktif. Selesaikan tahun akademik berjalan terlebih dahulu")
	ErrDeleteActiveYear       = errors.New("tahun akademik yang sedang aktif tidak dapat dihapus")
	ErrTargetYearNotAvailable = errors.New("tahun akademik baru belum tersedia untuk kenaikan kelas")
	ErrInvalidSemesterStatus  = errors.New("status semester tidak valid")
	ErrSemesterNotActive      = errors.New("semester tidak aktif, tidak dapat ditutup")
	ErrInvalidDate            = errors.New("tanggal akhir tidak boleh lebih awal dari tanggal mulai")
	ErrInvalidDateFormat      = errors.New("format tanggal tidak valid (Gunakan YYYY-MM-DD)")
)
