package class

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

// GetAll mengambil semua data kelas
// @Summary      Get all classes
// @Description  Get a list of all classes
// @Tags         Classes
// @Produce      json
// @Success      200 {object} response.DefaultResponse{data=[]ClassResponse}
// @Failure      500 {object} response.DefaultResponse
// @Router       /classes [get]
func (h *Handler) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	classes, err := h.service.GetAll(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil data kelas", classes)
}

// GetByID mengambil data kelas berdasarkan ID
// @Summary      Get class by ID
// @Description  Get detailed information of a class by its ID
// @Tags         Classes
// @Produce      json
// @Param        id   path      string  true  "Class ID"
// @Success      200 {object} response.DefaultResponse{data=ClassResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /classes/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ID kelas wajib diisi", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	class, err := h.service.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrClassNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil data kelas", class)
}

// Create menambahkan data kelas baru (Single)
// @Summary      Create new class
// @Description  Add a single new class to the system
// @Tags         Classes
// @Accept       json
// @Produce      json
// @Param        request body ClassRequest true "Class Data"
// @Success      201 {object} response.DefaultResponse{data=ClassResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      409 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /classes [post]
func (h *Handler) Create(c *gin.Context) {
	var req ClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.service.Create(ctx, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrClassNameExists) {
			statusCode = http.StatusConflict
		} else if errors.Is(err, ErrMajorRequired) {
			statusCode = http.StatusBadRequest
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Berhasil menambahkan kelas baru", res)
}

// CreateBatch membuat banyak kelas otomatis (Generative Alphabet)
// @Summary      Create classes in batch
// @Description  Auto-generate multiple classes alphabetically per level based on existing DB records
// @Tags         Classes
// @Accept       json
// @Produce      json
// @Param        request body ClassBatchRequest true "Batch Class Data"
// @Success      201 {object} response.DefaultResponse{data=[]ClassResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /classes/batch [post]
func (h *Handler) CreateBatch(c *gin.Context) {
	var req ClassBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second) // Batch butuh waktu lebih
	defer cancel()

	res, err := h.service.CreateBatch(ctx, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Berhasil membuat kelas secara massal", res)
}

// Update merubah data kelas (kuota / wali kelas)
// @Summary      Update class
// @Description  Update existing class capacity or homeroom teacher
// @Tags         Classes
// @Accept       json
// @Produce      json
// @Param        id      path     string             true  "Class ID"
// @Param        request body     ClassUpdateRequest true  "Class Data"
// @Success      200 {object} response.DefaultResponse{data=ClassResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /classes/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ID kelas wajib diisi", nil)
		return
	}

	var req ClassUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.service.Update(ctx, id, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrClassNotFound) {
			statusCode = http.StatusNotFound
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil memperbarui data kelas", res)
}

// Delete menghapus data kelas secara soft delete
// @Summary      Delete class
// @Description  Soft delete a class by ID
// @Tags         Classes
// @Produce      json
// @Param        id   path      string  true  "Class ID"
// @Success      200 {object} response.DefaultResponse
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /classes/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ID kelas wajib diisi", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err := h.service.Delete(ctx, id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrClassNotFound) {
			statusCode = http.StatusNotFound
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil menghapus data kelas", nil)
}
