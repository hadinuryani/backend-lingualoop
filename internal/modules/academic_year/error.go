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

	// Validasi bisnis tahun akademik
	ErrInvalidYearFormat     = errors.New("format tahun harus YYYY/YYYY (contoh: 2026/2027)")
	ErrInvalidYearSequence   = errors.New("tahun kedua harus = tahun pertama + 1 (contoh: 2026/2027)")
	ErrStartDateYearMismatch = errors.New("tahun start_date harus sesuai dengan tahun pertama di field year")
	ErrEndDateYearMismatch   = errors.New("tahun end_date harus sesuai dengan tahun kedua di field year")
	ErrAcademicRangeTooLong  = errors.New("rentang tahun akademik terlalu panjang (maksimal 13 bulan)")
	ErrSemesterOutOfRange    = errors.New("tanggal semester harus berada dalam rentang tahun akademik")
	ErrSemesterDateOrder     = errors.New("urutan tanggal semester harus: start → end_kbm → assessment")
	ErrOddBeforeEven         = errors.New("semester ganjil harus dimulai lebih dulu dari semester genap")
	ErrDraftExists           = errors.New("masih ada tahun akademik dalam status draft. harap selesaikan atau hapus terlebih dahulu")
)
