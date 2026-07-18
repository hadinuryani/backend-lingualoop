package region

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

type Repository interface {
	// Provinces
	FindAllProvinces(ctx context.Context, search string) ([]*Province, error)
	FindProvinceByID(ctx context.Context, id int) (*Province, error)

	// Cities
	FindCitiesByProvinceID(ctx context.Context, provinceID int, search string) ([]*CityResponse, error)
	FindCityByID(ctx context.Context, id int) (*CityResponse, error)

	// Districts
	FindDistrictsByCityID(ctx context.Context, cityID int, search string) ([]*DistrictResponse, error)
	FindDistrictByID(ctx context.Context, id int) (*DistrictResponse, error)

	// Subdistricts
	FindSubdistrictsByDistrictID(ctx context.Context, districtID int, search string) ([]*SubdistrictResponse, error)
	FindSubdistrictByID(ctx context.Context, id int) (*SubdistrictResponse, error)

	// Postal Codes
	FindPostalCodesBySubdistrictID(ctx context.Context, subdistrictID int) ([]*PostalCodeResponse, error)
	SearchByPostalCode(ctx context.Context, postalCode string) ([]*PostalCodeResponse, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}

// ======================== PROVINCES ========================

func (r *repository) FindAllProvinces(ctx context.Context, search string) ([]*Province, error) {
	query := `SELECT id, name FROM provinces`
	var args []interface{}

	if search != "" {
		query += ` WHERE name LIKE ?`
		args = append(args, "%"+strings.ToUpper(search)+"%")
	}

	query += ` ORDER BY name ASC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var provinces []*Province
	for rows.Next() {
		var p Province
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			return nil, err
		}
		provinces = append(provinces, &p)
	}
	return provinces, rows.Err()
}

func (r *repository) FindProvinceByID(ctx context.Context, id int) (*Province, error) {
	query := `SELECT id, name FROM provinces WHERE id = ?`
	var p Province
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProvinceNotFound
		}
		return nil, err
	}
	return &p, nil
}

// ======================== CITIES ========================

func (r *repository) FindCitiesByProvinceID(ctx context.Context, provinceID int, search string) ([]*CityResponse, error) {
	query := `
		SELECT c.id, c.name, c.province_id, p.name AS province_name
		FROM cities c
		JOIN provinces p ON c.province_id = p.id
		WHERE c.province_id = ?
	`
	args := []interface{}{provinceID}

	if search != "" {
		query += ` AND c.name LIKE ?`
		args = append(args, "%"+strings.ToUpper(search)+"%")
	}

	query += ` ORDER BY c.name ASC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []*CityResponse
	for rows.Next() {
		var c CityResponse
		if err := rows.Scan(&c.ID, &c.Name, &c.ProvinceID, &c.ProvinceName); err != nil {
			return nil, err
		}
		cities = append(cities, &c)
	}
	return cities, rows.Err()
}

func (r *repository) FindCityByID(ctx context.Context, id int) (*CityResponse, error) {
	query := `
		SELECT c.id, c.name, c.province_id, p.name AS province_name
		FROM cities c
		JOIN provinces p ON c.province_id = p.id
		WHERE c.id = ?
	`
	var c CityResponse
	err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.Name, &c.ProvinceID, &c.ProvinceName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCityNotFound
		}
		return nil, err
	}
	return &c, nil
}

// ======================== DISTRICTS ========================

func (r *repository) FindDistrictsByCityID(ctx context.Context, cityID int, search string) ([]*DistrictResponse, error) {
	query := `
		SELECT d.id, d.name, d.city_id, c.name AS city_name
		FROM districts d
		JOIN cities c ON d.city_id = c.id
		WHERE d.city_id = ?
	`
	args := []interface{}{cityID}

	if search != "" {
		query += ` AND d.name LIKE ?`
		args = append(args, "%"+strings.ToUpper(search)+"%")
	}

	query += ` ORDER BY d.name ASC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var districts []*DistrictResponse
	for rows.Next() {
		var d DistrictResponse
		if err := rows.Scan(&d.ID, &d.Name, &d.CityID, &d.CityName); err != nil {
			return nil, err
		}
		districts = append(districts, &d)
	}
	return districts, rows.Err()
}

func (r *repository) FindDistrictByID(ctx context.Context, id int) (*DistrictResponse, error) {
	query := `
		SELECT d.id, d.name, d.city_id, c.name AS city_name
		FROM districts d
		JOIN cities c ON d.city_id = c.id
		WHERE d.id = ?
	`
	var d DistrictResponse
	err := r.db.QueryRowContext(ctx, query, id).Scan(&d.ID, &d.Name, &d.CityID, &d.CityName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDistrictNotFound
		}
		return nil, err
	}
	return &d, nil
}

// ======================== SUBDISTRICTS ========================

func (r *repository) FindSubdistrictsByDistrictID(ctx context.Context, districtID int, search string) ([]*SubdistrictResponse, error) {
	query := `
		SELECT s.id, s.name, s.district_id, d.name AS district_name
		FROM subdistricts s
		JOIN districts d ON s.district_id = d.id
		WHERE s.district_id = ?
	`
	args := []interface{}{districtID}

	if search != "" {
		query += ` AND s.name LIKE ?`
		args = append(args, "%"+strings.ToUpper(search)+"%")
	}

	query += ` ORDER BY s.name ASC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subdistricts []*SubdistrictResponse
	for rows.Next() {
		var s SubdistrictResponse
		if err := rows.Scan(&s.ID, &s.Name, &s.DistrictID, &s.DistrictName); err != nil {
			return nil, err
		}
		subdistricts = append(subdistricts, &s)
	}
	return subdistricts, rows.Err()
}

func (r *repository) FindSubdistrictByID(ctx context.Context, id int) (*SubdistrictResponse, error) {
	query := `
		SELECT s.id, s.name, s.district_id, d.name AS district_name
		FROM subdistricts s
		JOIN districts d ON s.district_id = d.id
		WHERE s.id = ?
	`
	var s SubdistrictResponse
	err := r.db.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.Name, &s.DistrictID, &s.DistrictName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSubdistrictNotFound
		}
		return nil, err
	}
	return &s, nil
}

// ======================== POSTAL CODES ========================

func (r *repository) FindPostalCodesBySubdistrictID(ctx context.Context, subdistrictID int) ([]*PostalCodeResponse, error) {
	query := `
		SELECT 
			pc.id, pc.subdistrict_id, s.name AS subdistrict_name,
			pc.district_id, d.name AS district_name,
			pc.city_id, c.name AS city_name,
			pc.province_id, p.name AS province_name,
			pc.postal_code
		FROM postal_codes pc
		JOIN subdistricts s ON pc.subdistrict_id = s.id
		JOIN districts d ON pc.district_id = d.id
		JOIN cities c ON pc.city_id = c.id
		JOIN provinces p ON pc.province_id = p.id
		WHERE pc.subdistrict_id = ?
		ORDER BY pc.postal_code ASC
	`

	rows, err := r.db.QueryContext(ctx, query, subdistrictID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var postalCodes []*PostalCodeResponse
	for rows.Next() {
		var pc PostalCodeResponse
		if err := rows.Scan(
			&pc.ID, &pc.SubdistrictID, &pc.SubdistrictName,
			&pc.DistrictID, &pc.DistrictName,
			&pc.CityID, &pc.CityName,
			&pc.ProvinceID, &pc.ProvinceName,
			&pc.PostalCode,
		); err != nil {
			return nil, err
		}
		postalCodes = append(postalCodes, &pc)
	}
	return postalCodes, rows.Err()
}

func (r *repository) SearchByPostalCode(ctx context.Context, postalCode string) ([]*PostalCodeResponse, error) {
	if len(postalCode) < 3 {
		return nil, ErrInvalidPostalCode
	}

	query := `
		SELECT 
			pc.id, pc.subdistrict_id, s.name AS subdistrict_name,
			pc.district_id, d.name AS district_name,
			pc.city_id, c.name AS city_name,
			pc.province_id, p.name AS province_name,
			pc.postal_code
		FROM postal_codes pc
		JOIN subdistricts s ON pc.subdistrict_id = s.id
		JOIN districts d ON pc.district_id = d.id
		JOIN cities c ON pc.city_id = c.id
		JOIN provinces p ON pc.province_id = p.id
		WHERE pc.postal_code LIKE ?
		ORDER BY pc.postal_code ASC, p.name ASC, c.name ASC
		LIMIT 100
	`

	rows, err := r.db.QueryContext(ctx, query, postalCode+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var postalCodes []*PostalCodeResponse
	for rows.Next() {
		var pc PostalCodeResponse
		if err := rows.Scan(
			&pc.ID, &pc.SubdistrictID, &pc.SubdistrictName,
			&pc.DistrictID, &pc.DistrictName,
			&pc.CityID, &pc.CityName,
			&pc.ProvinceID, &pc.ProvinceName,
			&pc.PostalCode,
		); err != nil {
			return nil, err
		}
		postalCodes = append(postalCodes, &pc)
	}
	return postalCodes, rows.Err()
}
