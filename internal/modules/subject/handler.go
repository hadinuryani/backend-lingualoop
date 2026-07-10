package subject

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

// GetAll mengambil semua data mata pelajaran
// @Summary      Get all subjects
// @Description  Get a list of all subjects
// @Tags         Subjects
// @Produce      json
// @Success      200 {object} response.DefaultResponse{data=[]SubjectResponse}
// @Failure      500 {object} response.DefaultResponse
// @Router       /subjects [get]
func (h *Handler) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	subjects, err := h.service.GetAll(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil data mata pelajaran", subjects)
}

// GetByID mengambil data mata pelajaran berdasarkan ID
// @Summary      Get subject by ID
// @Description  Get detailed information of a subject by its ID
// @Tags         Subjects
// @Produce      json
// @Param        id   path      string  true  "Subject ID"
// @Success      200 {object} response.DefaultResponse{data=SubjectResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /subjects/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ID mata pelajaran wajib diisi", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	sub, err := h.service.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrSubjectNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil data mata pelajaran", sub)
}

// Create menambahkan data mata pelajaran baru
// @Summary      Create new subject
// @Description  Add a new subject to the system
// @Tags         Subjects
// @Accept       json
// @Produce      json
// @Param        request body SubjectRequest true "Subject Data"
// @Success      201 {object} response.DefaultResponse{data=SubjectResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      409 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /subjects [post]
func (h *Handler) Create(c *gin.Context) {
	var req SubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.service.Create(ctx, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrSubjectCodeExists) {
			statusCode = http.StatusConflict
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Berhasil menambahkan mata pelajaran baru", res)
}

// Update merubah data mata pelajaran
// @Summary      Update subject
// @Description  Update existing subject data
// @Tags         Subjects
// @Accept       json
// @Produce      json
// @Param        id      path     string         true  "Subject ID"
// @Param        request body     SubjectRequest true  "Subject Data"
// @Success      200 {object} response.DefaultResponse{data=SubjectResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      409 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /subjects/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ID mata pelajaran wajib diisi", nil)
		return
	}

	var req SubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.service.Update(ctx, id, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrSubjectNotFound) {
			statusCode = http.StatusNotFound
		} else if errors.Is(err, ErrSubjectCodeExists) {
			statusCode = http.StatusConflict
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil memperbarui data mata pelajaran", res)
}

// Delete menghapus data mata pelajaran secara soft delete
// @Summary      Delete subject
// @Description  Soft delete a subject by ID
// @Tags         Subjects
// @Produce      json
// @Param        id   path      string  true  "Subject ID"
// @Success      200 {object} response.DefaultResponse
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /subjects/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ID mata pelajaran wajib diisi", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err := h.service.Delete(ctx, id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrSubjectNotFound) {
			statusCode = http.StatusNotFound
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil menghapus data mata pelajaran", nil)
}
