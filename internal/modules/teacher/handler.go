package teacher

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

// GetAll mengambil semua data guru
// @Summary      Get all teachers
// @Description  Get a list of all teachers
// @Tags         Teachers
// @Produce      json
// @Success      200 {object} response.DefaultResponse{data=[]TeacherResponse}
// @Failure      500 {object} response.DefaultResponse
// @Router       /teachers [get]
func (h *Handler) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	teachers, err := h.service.GetAll(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil data guru", teachers)
}

// GetByID mengambil data guru berdasarkan ID
// @Summary      Get teacher by ID
// @Description  Get detailed information of a teacher by their ID
// @Tags         Teachers
// @Produce      json
// @Param        id   path      string  true  "Teacher ID"
// @Success      200 {object} response.DefaultResponse{data=TeacherResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /teachers/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ID guru wajib diisi", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	teacher, err := h.service.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrTeacherNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil data guru", teacher)
}

// Create menambahkan data guru baru
// @Summary      Create new teacher
// @Description  Add a new teacher to the system and auto-generate user account
// @Tags         Teachers
// @Accept       json
// @Produce      json
// @Param        request body TeacherRequest true "Teacher Data"
// @Success      201 {object} response.DefaultResponse{data=TeacherResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      409 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /teachers [post]
func (h *Handler) Create(c *gin.Context) {
	var req TeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.service.Create(ctx, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrNipExists) {
			statusCode = http.StatusConflict
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Berhasil menambahkan guru baru", res)
}

// Update merubah data guru
// @Summary      Update teacher
// @Description  Update existing teacher information
// @Tags         Teachers
// @Accept       json
// @Produce      json
// @Param        id      path     string          true  "Teacher ID"
// @Param        request body     TeacherRequest  true  "Teacher Data"
// @Success      200 {object} response.DefaultResponse{data=TeacherResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      409 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /teachers/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ID guru wajib diisi", nil)
		return
	}

	var req TeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.service.Update(ctx, id, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrTeacherNotFound) {
			statusCode = http.StatusNotFound
		} else if errors.Is(err, ErrNipExists) {
			statusCode = http.StatusConflict
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil memperbarui data guru", res)
}

// ToggleStatus mengaktifkan/menonaktifkan guru
// @Summary      Toggle teacher status
// @Description  Change teacher status between ACTIVE and INACTIVE
// @Tags         Teachers
// @Produce      json
// @Param        id   path      string  true  "Teacher ID"
// @Success      200 {object} response.DefaultResponse{data=TeacherResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /teachers/{id}/status [patch]
func (h *Handler) ToggleStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ID guru wajib diisi", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.service.ToggleStatus(ctx, id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrTeacherNotFound) {
			statusCode = http.StatusNotFound
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil merubah status guru", res)
}

// Delete menghapus data guru secara soft delete
// @Summary      Delete teacher
// @Description  Soft delete a teacher by ID
// @Tags         Teachers
// @Produce      json
// @Param        id   path      string  true  "Teacher ID"
// @Success      200 {object} response.DefaultResponse
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /teachers/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ID guru wajib diisi", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err := h.service.Delete(ctx, id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrTeacherNotFound) {
			statusCode = http.StatusNotFound
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil menghapus data guru", nil)
}
