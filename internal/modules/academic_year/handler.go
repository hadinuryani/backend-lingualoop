package academic_year

import (
	"context"
	"errors"
	"net/http"
	"time"

	"backend-lingualoop/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service         Service
	templateService TemplateService
}

func NewHandler(service Service, templateService TemplateService) *Handler {
	return &Handler{
		service:         service,
		templateService: templateService,
	}
}

// isValidationError mengecek apakah error berasal dari validasi bisnis (bukan error sistem).
// Error validasi harus dikembalikan sebagai 400 Bad Request, bukan 500.
func isValidationError(err error) bool {
	validationErrors := []error{
		ErrInvalidDate,
		ErrInvalidDateFormat,
		ErrInvalidYearFormat,
		ErrInvalidYearSequence,
		ErrStartDateYearMismatch,
		ErrEndDateYearMismatch,
		ErrAcademicRangeTooLong,
		ErrSemesterOutOfRange,
		ErrSemesterDateOrder,
		ErrOddBeforeEven,
		ErrInvalidSemesterStatus,
		ErrSemesterNotActive,
	}
	for _, ve := range validationErrors {
		if errors.Is(err, ve) {
			return true
		}
	}
	return false
}

// GetAll mengambil semua data tahun akademik
// @Summary      Get all academic years
// @Description  Get a list of all academic years including nested semester data
// @Tags         Academic Years
// @Produce      json
// @Success      200 {object} response.DefaultResponse{data=[]AcademicYearResponse}
// @Failure      500 {object} response.DefaultResponse
// @Router       /academic-years [get]
// @Security      BearerAuth
func (h *Handler) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	years, err := h.service.GetAll(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil data tahun akademik", years)
}

// GetByID mengambil data tahun akademik berdasarkan ID
// @Summary      Get academic year by ID
// @Description  Get detailed information of an academic year by its ID
// @Tags         Academic Years
// @Produce      json
// @Param        id   path      string  true  "Academic Year ID"
// @Success      200 {object} response.DefaultResponse{data=AcademicYearResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /academic-years/{id} [get]
// @Security      BearerAuth
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ID tahun akademik wajib diisi", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	y, err := h.service.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrAcademicYearNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil data tahun akademik", y)
}


// Create menambahkan data tahun akademik baru
// @Summary      Create new academic year
// @Description  Add a new academic year to the system
// @Tags         Academic Years
// @Accept       json
// @Produce      json
// @Param        request body AcademicYearRequest true "Academic Year Data"
// @Success      201 {object} response.DefaultResponse{data=AcademicYearResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      409 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /academic-years [post]
// @Security      BearerAuth
func (h *Handler) Create(c *gin.Context) {
	var req AcademicYearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.service.Create(ctx, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrAcademicYearExists) {
			statusCode = http.StatusConflict
		} else if errors.Is(err, ErrDraftExists) {
			statusCode = http.StatusConflict
		} else if isValidationError(err) {
			statusCode = http.StatusBadRequest
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Berhasil menambahkan tahun akademik baru", res)
}

// Update merubah data tahun akademik
// @Summary      Update academic year
// @Description  Update existing academic year data
// @Tags         Academic Years
// @Accept       json
// @Produce      json
// @Param        id      path     string              true  "Academic Year ID"
// @Param        request body     AcademicYearRequest true  "Academic Year Data"
// @Success      200 {object} response.DefaultResponse{data=AcademicYearResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      409 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /academic-years/{id} [put]
// @Security      BearerAuth
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ID tahun akademik wajib diisi", nil)
		return
	}

	var req AcademicYearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.service.Update(ctx, id, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrAcademicYearNotFound) {
			statusCode = http.StatusNotFound
		} else if errors.Is(err, ErrAcademicYearExists) {
			statusCode = http.StatusConflict
		} else if isValidationError(err) {
			statusCode = http.StatusBadRequest
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil memperbarui data tahun akademik", res)
}

// Activate mengaktifkan tahun akademik
// @Summary      Activate academic year
// @Description  Activate an academic year and set semester ganjil to active. Fails if another year is already active.
// @Tags         Academic Years
// @Produce      json
// @Param        id   path      string  true  "Academic Year ID"
// @Success      200 {object} response.DefaultResponse{data=AcademicYearResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      409 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /academic-years/{id}/activate [put]
// @Security      BearerAuth
func (h *Handler) Activate(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ID tahun akademik wajib diisi", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.service.Activate(ctx, id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrAcademicYearNotFound) {
			statusCode = http.StatusNotFound
		} else if errors.Is(err, ErrMultipleActiveYears) {
			statusCode = http.StatusConflict
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengaktifkan tahun akademik", res)
}

// GetTemplate menghasilkan template tahun akademik baru berdasarkan tahun sebelumnya
// @Summary      Get academic year template
// @Description  Generate a template for the next academic year with dates shifted by +1 year. Fails if a draft already exists.
// @Tags         Academic Years
// @Produce      json
// @Success      200 {object} response.DefaultResponse{data=TemplateResponse}
// @Failure      409 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /academic-years/template [get]
// @Security     BearerAuth
func (h *Handler) GetTemplate(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.templateService.GenerateNext(ctx)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrDraftExists) {
			statusCode = http.StatusConflict
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil membuat template tahun akademik", res)
}

// UpdateSemesterStatus mengubah status sub-semester (ganjil/genap)
// @Summary      Update semester status
// @Description  Change the status of ganjil or genap semester
// @Tags         Academic Years
// @Accept       json
// @Produce      json
// @Param        id      path     string                true  "Academic Year ID"
// @Param        request body     SemesterStatusRequest true  "Semester Status Data"
// @Success      200 {object} response.DefaultResponse{data=AcademicYearResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /academic-years/{id}/semester-status [put]
// @Security      BearerAuth
func (h *Handler) UpdateSemesterStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ID tahun akademik wajib diisi", nil)
		return
	}

	var req SemesterStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.service.UpdateSemesterStatus(ctx, id, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrAcademicYearNotFound) {
			statusCode = http.StatusNotFound
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil memperbarui status semester", res)
}

// CloseSemester menutup semester (ganjil/genap) dan memicu state selanjutnya
// @Summary      Close semester
// @Description  Close a semester (if ganjil -> unlocks genap, if genap -> changes year status to Menunggu Kenaikan)
// @Tags         Academic Years
// @Accept       json
// @Produce      json
// @Param        id      path     string               true  "Academic Year ID"
// @Param        request body     CloseSemesterRequest true  "Close Semester Data"
// @Success      200 {object} response.DefaultResponse{data=AcademicYearResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /academic-years/{id}/close-semester [put]
// @Security      BearerAuth
func (h *Handler) CloseSemester(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ID tahun akademik wajib diisi", nil)
		return
	}

	var req CloseSemesterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.service.CloseSemester(ctx, id, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrAcademicYearNotFound) {
			statusCode = http.StatusNotFound
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil menutup semester", res)
}

// Delete menghapus data tahun akademik secara soft delete
// @Summary      Delete academic year
// @Description  Soft delete an academic year by ID. Fails if the year is active.
// @Tags         Academic Years
// @Produce      json
// @Param        id   path      string  true  "Academic Year ID"
// @Success      200 {object} response.DefaultResponse
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      409 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /academic-years/{id} [delete]
// @Security      BearerAuth
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ID tahun akademik wajib diisi", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err := h.service.Delete(ctx, id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrAcademicYearNotFound) {
			statusCode = http.StatusNotFound
		} else if errors.Is(err, ErrDeleteActiveYear) {
			statusCode = http.StatusConflict
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil menghapus data tahun akademik", nil)
}

func (h *Handler) FinalizePromotion(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "academic year ID is required", nil)
		return
	}

	var req FinalizePromotionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request format", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second) // Waktu dinaikin karena ini operasi batch yang berat
	defer cancel()

	if err := h.service.FinalizePromotion(ctx, id, req); err != nil {
		if errors.Is(err, ErrAcademicYearNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if err.Error() == "academic year is not ready for promotion" || err.Error() == "target academic year not found" {
			response.Error(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		if err.Error() == "promotion already processed or is running" {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to finalize promotion", nil)
		return
	}

	response.Success(c, http.StatusOK, "promotion process completed successfully", nil)
}
