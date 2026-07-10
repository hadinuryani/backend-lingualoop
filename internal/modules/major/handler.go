package major

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
	return &Handler{service}
}

// GetAll godoc
// @Summary      Get All Majors
// @Description  Retrieve a list of all majors
// @Tags         Majors
// @Produce      json
// @Success      200 {object} response.DefaultResponse{data=[]MajorResponse}
// @Failure      500 {object} response.DefaultResponse
// @Router       /majors [get]
func (h *Handler) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	majors, err := h.service.GetAll(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, response.MsgSuccess, majors)
}

// Create godoc
// @Summary      Create Major
// @Description  Create a new major
// @Tags         Majors
// @Accept       json
// @Produce      json
// @Param        request body MajorRequest true "Major Data"
// @Success      201 {object} response.DefaultResponse{data=MajorResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      409 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /majors [post]
func (h *Handler) Create(c *gin.Context) {
	var req MajorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.service.Create(ctx, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrCodeExists) || errors.Is(err, ErrNameExists) {
			statusCode = http.StatusConflict // 409 Conflict
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, response.MsgCreated, res)
}

// Update godoc
// @Summary      Update Major
// @Description  Update an existing major by ID
// @Tags         Majors
// @Accept       json
// @Produce      json
// @Param        id path string true "Major ID"
// @Param        request body MajorRequest true "Major Data"
// @Success      200 {object} response.DefaultResponse{data=MajorResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      409 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /majors/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ID jurusan wajib diisi", nil)
		return
	}

	var req MajorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.service.Update(ctx, id, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrCodeExists) || errors.Is(err, ErrNameExists) {
			statusCode = http.StatusConflict // 409 Conflict
		} else if errors.Is(err, ErrMajorNotFound) {
			statusCode = http.StatusNotFound // 404 Not Found
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, response.MsgUpdated, res)
}

// Delete godoc
// @Summary      Delete Major
// @Description  Delete a major by ID
// @Tags         Majors
// @Produce      json
// @Param        id path string true "Major ID"
// @Success      200 {object} response.DefaultResponse
// @Failure      400 {object} response.DefaultResponse
// @Failure      404 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /majors/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ID jurusan wajib diisi", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err := h.service.Delete(ctx, id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrMajorNotFound) {
			statusCode = http.StatusNotFound
		}
		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, response.MsgDeleted, nil)
}
