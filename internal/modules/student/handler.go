package student

import (
	"context"
	"errors"
	"net/http"
	"time"

	"backend-lingualoop/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetAll mengambil semua data siswa
// @Summary      Get all students
// @Description  Get a list of all students
// @Tags         Students
// @Produce      json
// @Success      200 {object} response.DefaultResponse{data=[]StudentResponse}
// @Failure      500 {object} response.DefaultResponse
// @Router       /students [get]
// @Security      BearerAuth
func (h *Handler) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	students, err := h.service.GetAll(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil data siswa", students)
}

// GetByID mengambil data siswa berdasarkan ID
// @Summary      Get student by ID
// @Description  Get detailed information of a student by their ID
// @Tags         Students
// @Produce      json
// @Param        id   path      string  true  "Student ID"
// @Success      200 {object} response.DefaultResponse{data=StudentResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /students/{id} [get]
// @Security      BearerAuth
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ID siswa wajib diisi", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	student, err := h.service.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil data siswa", student)
}

// Create menambahkan data siswa baru
// @Summary      Create new student
// @Description  Add a new student to the system and auto-generate user account
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        request body StudentRequest true "Student Data"
// @Success      201 {object} response.DefaultResponse{data=StudentResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      409 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /students [post]
// @Security      BearerAuth
func (h *Handler) Create(c *gin.Context) {
	var req StudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.service.Create(ctx, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrNisExists) {
			statusCode = http.StatusConflict
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Berhasil menambahkan siswa baru", res)
}

// Update merubah data siswa
// @Summary      Update student
// @Description  Update existing student information
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        id      path     string          true  "Student ID"
// @Param        request body     StudentRequest  true  "Student Data"
// @Success      200 {object} response.DefaultResponse{data=StudentResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      409 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /students/{id} [put]
// @Security      BearerAuth
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ID siswa wajib diisi", nil)
		return
	}

	var req StudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.service.Update(ctx, id, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrStudentNotFound) {
			statusCode = http.StatusNotFound
		} else if errors.Is(err, ErrNisExists) {
			statusCode = http.StatusConflict
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil memperbarui data siswa", res)
}

// UpdateStatus mengubah status siswa (ACTIVE/GRADUATED/TRANSFER/INACTIVE)
// @Summary      Update student status
// @Description  Change student status and auto-sync login capabilities
// @Tags         Students
// @Produce      json
// @Param        id      path     string  true  "Student ID"
// @Param        status  query    string  true  "New Status (ACTIVE/GRADUATED/TRANSFER/INACTIVE)"
// @Success      200 {object} response.DefaultResponse{data=StudentResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /students/{id}/status [patch]
// @Security      BearerAuth
func (h *Handler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ID siswa wajib diisi", nil)
		return
	}

	// Format untuk PATCH di frontend biasanya via body `{ status: "ACTIVE" }`
	// Tapi kita ikuti pola sederhana
	var payload struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.ValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.service.UpdateStatus(ctx, id, payload.Status)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrStudentNotFound) {
			statusCode = http.StatusNotFound
		} else if errors.Is(err, ErrInvalidStatus) {
			statusCode = http.StatusBadRequest
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil merubah status siswa", res)
}

// Delete menghapus data siswa secara soft delete
// @Summary      Delete student
// @Description  Soft delete a student by ID
// @Tags         Students
// @Produce      json
// @Param        id   path      string  true  "Student ID"
// @Success      200 {object} response.DefaultResponse
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /students/{id} [delete]
// @Security      BearerAuth
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ID siswa wajib diisi", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err := h.service.Delete(ctx, id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrStudentNotFound) {
			statusCode = http.StatusNotFound
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil menghapus data siswa", nil)
}
