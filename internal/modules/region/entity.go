package region

// Province adalah representasi Entity dari tabel `provinces` di database.
type Province struct {
	ID   int
	Name string
}

// City adalah representasi Entity dari tabel `cities` di database.
type City struct {
	ID         int
	Name       string
	ProvinceID int
}

// District adalah representasi Entity dari tabel `districts` di database.
type District struct {
	ID     int
	Name   string
	CityID int
}

// Subdistrict adalah representasi Entity dari tabel `subdistricts` di database.
type Subdistrict struct {
	ID         int
	Name       string
	DistrictID int
}

// PostalCode adalah representasi Entity dari tabel `postal_codes` di database.
type PostalCode struct {
	ID            int
	SubdistrictID int
	DistrictID    int
	CityID        int
	ProvinceID    int
	PostalCode    string
}
