package region

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"backend-lingualoop/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetAllProvinces godoc
// @Summary      Get All Provinces
// @Description  Retrieve a list of all provinces. Supports optional search by name.
// @Tags         Region
// @Produce      json
// @Param        search query string false "Search province by name"
// @Success      200 {object} response.DefaultResponse{data=[]ProvinceResponse}
// @Failure      500 {object} response.DefaultResponse
// @Router       /region/provinces [get]
// @Security      BearerAuth
func (h *Handler) GetAllProvinces(c *gin.Context) {
	search := c.Query("search")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	provinces, err := h.service.GetAllProvinces(ctx, search)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, response.MsgSuccess, provinces)
}

// GetProvinceByID godoc
// @Summary      Get Province by ID
// @Description  Retrieve province detail by ID
// @Tags         Region
// @Produce      json
// @Param        id path int true "Province ID"
// @Success      200 {object} response.DefaultResponse{data=ProvinceResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /region/provinces/{id} [get]
// @Security      BearerAuth
func (h *Handler) GetProvinceByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID provinsi tidak valid", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	province, err := h.service.GetProvinceByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrProvinceNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, response.MsgSuccess, province)
}

// GetCitiesByProvinceID godoc
// @Summary      Get Cities by Province ID
// @Description  Retrieve all cities/regencies within a province. Supports optional search by name.
// @Tags         Region
// @Produce      json
// @Param        id path int true "Province ID"
// @Param        search query string false "Search city by name"
// @Success      200 {object} response.DefaultResponse{data=[]CityResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /region/provinces/{id}/cities [get]
// @Security      BearerAuth
func (h *Handler) GetCitiesByProvinceID(c *gin.Context) {
	provinceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID provinsi tidak valid", nil)
		return
	}

	search := c.Query("search")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	cities, err := h.service.GetCitiesByProvinceID(ctx, provinceID, search)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, response.MsgSuccess, cities)
}

// GetCityByID godoc
// @Summary      Get City by ID
// @Description  Retrieve city detail by ID
// @Tags         Region
// @Produce      json
// @Param        id path int true "City ID"
// @Success      200 {object} response.DefaultResponse{data=CityResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /region/cities/{id} [get]
// @Security      BearerAuth
func (h *Handler) GetCityByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID kota tidak valid", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	city, err := h.service.GetCityByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrCityNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, response.MsgSuccess, city)
}

// GetDistrictsByCityID godoc
// @Summary      Get Districts by City ID
// @Description  Retrieve all districts within a city. Supports optional search by name.
// @Tags         Region
// @Produce      json
// @Param        id path int true "City ID"
// @Param        search query string false "Search district by name"
// @Success      200 {object} response.DefaultResponse{data=[]DistrictResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /region/cities/{id}/districts [get]
// @Security      BearerAuth
func (h *Handler) GetDistrictsByCityID(c *gin.Context) {
	cityID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID kota tidak valid", nil)
		return
	}

	search := c.Query("search")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	districts, err := h.service.GetDistrictsByCityID(ctx, cityID, search)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, response.MsgSuccess, districts)
}

// GetDistrictByID godoc
// @Summary      Get District by ID
// @Description  Retrieve district detail by ID
// @Tags         Region
// @Produce      json
// @Param        id path int true "District ID"
// @Success      200 {object} response.DefaultResponse{data=DistrictResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /region/districts/{id} [get]
// @Security      BearerAuth
func (h *Handler) GetDistrictByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID kecamatan tidak valid", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	district, err := h.service.GetDistrictByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrDistrictNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, response.MsgSuccess, district)
}

// GetSubdistrictsByDistrictID godoc
// @Summary      Get Subdistricts by District ID
// @Description  Retrieve all subdistricts/villages within a district. Supports optional search by name.
// @Tags         Region
// @Produce      json
// @Param        id path int true "District ID"
// @Param        search query string false "Search subdistrict by name"
// @Success      200 {object} response.DefaultResponse{data=[]SubdistrictResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /region/districts/{id}/subdistricts [get]
// @Security      BearerAuth
func (h *Handler) GetSubdistrictsByDistrictID(c *gin.Context) {
	districtID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID kecamatan tidak valid", nil)
		return
	}

	search := c.Query("search")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	subdistricts, err := h.service.GetSubdistrictsByDistrictID(ctx, districtID, search)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, response.MsgSuccess, subdistricts)
}

// GetSubdistrictByID godoc
// @Summary      Get Subdistrict by ID
// @Description  Retrieve subdistrict detail by ID
// @Tags         Region
// @Produce      json
// @Param        id path int true "Subdistrict ID"
// @Success      200 {object} response.DefaultResponse{data=SubdistrictResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /region/subdistricts/{id} [get]
// @Security      BearerAuth
func (h *Handler) GetSubdistrictByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID kelurahan tidak valid", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	subdistrict, err := h.service.GetSubdistrictByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrSubdistrictNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, response.MsgSuccess, subdistrict)
}

// GetPostalCodesBySubdistrictID godoc
// @Summary      Get Postal Codes by Subdistrict ID
// @Description  Retrieve postal codes for a subdistrict
// @Tags         Region
// @Produce      json
// @Param        id path int true "Subdistrict ID"
// @Success      200 {object} response.DefaultResponse{data=[]PostalCodeResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /region/subdistricts/{id}/postal-codes [get]
// @Security      BearerAuth
func (h *Handler) GetPostalCodesBySubdistrictID(c *gin.Context) {
	subdistrictID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID kelurahan tidak valid", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	postalCodes, err := h.service.GetPostalCodesBySubdistrictID(ctx, subdistrictID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, response.MsgSuccess, postalCodes)
}

// SearchPostalCode godoc
// @Summary      Search by Postal Code
// @Description  Search locations by postal code (minimum 3 characters)
// @Tags         Region
// @Produce      json
// @Param        code query string true "Postal code to search (min 3 chars)"
// @Success      200 {object} response.DefaultResponse{data=[]PostalCodeResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /region/postal-codes/search [get]
// @Security      BearerAuth
func (h *Handler) SearchPostalCode(c *gin.Context) {
	code := c.Query("code")
	if len(code) < 3 {
		response.Error(c, http.StatusBadRequest, "Kode pos minimal 3 karakter", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	postalCodes, err := h.service.SearchByPostalCode(ctx, code)
	if err != nil {
		if errors.Is(err, ErrInvalidPostalCode) {
			response.Error(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, response.MsgSuccess, postalCodes)
}
