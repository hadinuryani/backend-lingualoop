package region

import "errors"

var (
	ErrSystemFail          = errors.New("terjadi kesalahan pada sistem, silakan coba lagi")
	ErrProvinceNotFound    = errors.New("data provinsi tidak ditemukan")
	ErrCityNotFound        = errors.New("data kota/kabupaten tidak ditemukan")
	ErrDistrictNotFound    = errors.New("data kecamatan tidak ditemukan")
	ErrSubdistrictNotFound = errors.New("data kelurahan/desa tidak ditemukan")
	ErrInvalidPostalCode   = errors.New("kode pos minimal 3 karakter")
)
