package region

// ProvinceResponse adalah DTO response untuk data provinsi.
type ProvinceResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CityResponse adalah DTO response untuk data kota/kabupaten.
type CityResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ProvinceID   int    `json:"province_id"`
	ProvinceName string `json:"province_name,omitempty"`
}

// DistrictResponse adalah DTO response untuk data kecamatan.
type DistrictResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	CityID   int    `json:"city_id"`
	CityName string `json:"city_name,omitempty"`
}

// SubdistrictResponse adalah DTO response untuk data kelurahan/desa.
type SubdistrictResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	DistrictID   int    `json:"district_id"`
	DistrictName string `json:"district_name,omitempty"`
}

// PostalCodeResponse adalah DTO response untuk data kode pos.
type PostalCodeResponse struct {
	ID              int    `json:"id"`
	SubdistrictID   int    `json:"subdistrict_id"`
	SubdistrictName string `json:"subdistrict_name,omitempty"`
	DistrictID      int    `json:"district_id"`
	DistrictName    string `json:"district_name,omitempty"`
	CityID          int    `json:"city_id"`
	CityName        string `json:"city_name,omitempty"`
	ProvinceID      int    `json:"province_id"`
	ProvinceName    string `json:"province_name,omitempty"`
	PostalCode      string `json:"postal_code"`
}
