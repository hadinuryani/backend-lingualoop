package region

import (
	"context"
	"errors"
	"log/slog"
)

type Service interface {
	// Provinces
	GetAllProvinces(ctx context.Context, search string) ([]*ProvinceResponse, error)
	GetProvinceByID(ctx context.Context, id int) (*ProvinceResponse, error)

	// Cities
	GetCitiesByProvinceID(ctx context.Context, provinceID int, search string) ([]*CityResponse, error)
	GetCityByID(ctx context.Context, id int) (*CityResponse, error)

	// Districts
	GetDistrictsByCityID(ctx context.Context, cityID int, search string) ([]*DistrictResponse, error)
	GetDistrictByID(ctx context.Context, id int) (*DistrictResponse, error)

	// Subdistricts
	GetSubdistrictsByDistrictID(ctx context.Context, districtID int, search string) ([]*SubdistrictResponse, error)
	GetSubdistrictByID(ctx context.Context, id int) (*SubdistrictResponse, error)

	// Postal Codes
	GetPostalCodesBySubdistrictID(ctx context.Context, subdistrictID int) ([]*PostalCodeResponse, error)
	SearchByPostalCode(ctx context.Context, postalCode string) ([]*PostalCodeResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// ======================== PROVINCES ========================

func (s *service) GetAllProvinces(ctx context.Context, search string) ([]*ProvinceResponse, error) {
	provinces, err := s.repo.FindAllProvinces(ctx, search)
	if err != nil {
		return nil, logSystemError("Failed to query provinces", err, "search", search)
	}

	var responses []*ProvinceResponse
	for _, p := range provinces {
		responses = append(responses, mapProvinceToDTO(p))
	}

	if responses == nil {
		responses = []*ProvinceResponse{}
	}

	return responses, nil
}

func (s *service) GetProvinceByID(ctx context.Context, id int) (*ProvinceResponse, error) {
	p, err := s.repo.FindProvinceByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrProvinceNotFound) {
			return nil, ErrProvinceNotFound
		}
		return nil, logSystemError("Failed to find province by ID", err, "id", id)
	}

	return mapProvinceToDTO(p), nil
}

// ======================== CITIES ========================

func (s *service) GetCitiesByProvinceID(ctx context.Context, provinceID int, search string) ([]*CityResponse, error) {
	cities, err := s.repo.FindCitiesByProvinceID(ctx, provinceID, search)
	if err != nil {
		return nil, logSystemError("Failed to query cities", err, "province_id", provinceID, "search", search)
	}

	if cities == nil {
		cities = []*CityResponse{}
	}

	return cities, nil
}

func (s *service) GetCityByID(ctx context.Context, id int) (*CityResponse, error) {
	city, err := s.repo.FindCityByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrCityNotFound) {
			return nil, ErrCityNotFound
		}
		return nil, logSystemError("Failed to find city by ID", err, "id", id)
	}
	return city, nil
}

// ======================== DISTRICTS ========================

func (s *service) GetDistrictsByCityID(ctx context.Context, cityID int, search string) ([]*DistrictResponse, error) {
	districts, err := s.repo.FindDistrictsByCityID(ctx, cityID, search)
	if err != nil {
		return nil, logSystemError("Failed to query districts", err, "city_id", cityID, "search", search)
	}

	if districts == nil {
		districts = []*DistrictResponse{}
	}

	return districts, nil
}

func (s *service) GetDistrictByID(ctx context.Context, id int) (*DistrictResponse, error) {
	district, err := s.repo.FindDistrictByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrDistrictNotFound) {
			return nil, ErrDistrictNotFound
		}
		return nil, logSystemError("Failed to find district by ID", err, "id", id)
	}
	return district, nil
}

// ======================== SUBDISTRICTS ========================

func (s *service) GetSubdistrictsByDistrictID(ctx context.Context, districtID int, search string) ([]*SubdistrictResponse, error) {
	subdistricts, err := s.repo.FindSubdistrictsByDistrictID(ctx, districtID, search)
	if err != nil {
		return nil, logSystemError("Failed to query subdistricts", err, "district_id", districtID, "search", search)
	}

	if subdistricts == nil {
		subdistricts = []*SubdistrictResponse{}
	}

	return subdistricts, nil
}

func (s *service) GetSubdistrictByID(ctx context.Context, id int) (*SubdistrictResponse, error) {
	subdistrict, err := s.repo.FindSubdistrictByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrSubdistrictNotFound) {
			return nil, ErrSubdistrictNotFound
		}
		return nil, logSystemError("Failed to find subdistrict by ID", err, "id", id)
	}
	return subdistrict, nil
}

// ======================== POSTAL CODES ========================

func (s *service) GetPostalCodesBySubdistrictID(ctx context.Context, subdistrictID int) ([]*PostalCodeResponse, error) {
	postalCodes, err := s.repo.FindPostalCodesBySubdistrictID(ctx, subdistrictID)
	if err != nil {
		return nil, logSystemError("Failed to query postal codes", err, "subdistrict_id", subdistrictID)
	}

	if postalCodes == nil {
		postalCodes = []*PostalCodeResponse{}
	}

	return postalCodes, nil
}

func (s *service) SearchByPostalCode(ctx context.Context, postalCode string) ([]*PostalCodeResponse, error) {
	postalCodes, err := s.repo.SearchByPostalCode(ctx, postalCode)
	if err != nil {
		if errors.Is(err, ErrInvalidPostalCode) {
			return nil, ErrInvalidPostalCode
		}
		return nil, logSystemError("Failed to search postal codes", err, "postal_code", postalCode)
	}

	if postalCodes == nil {
		postalCodes = []*PostalCodeResponse{}
	}

	return postalCodes, nil
}

// ======================== HELPERS ========================

func logSystemError(msg string, err error, attrs ...any) error {
	slog.Error(msg, append([]any{"error", err}, attrs...)...)
	return ErrSystemFail
}

func mapProvinceToDTO(p *Province) *ProvinceResponse {
	if p == nil {
		return nil
	}
	return &ProvinceResponse{
		ID:   p.ID,
		Name: p.Name,
	}
}
