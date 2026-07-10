package auth

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

// Login godoc
// @Summary      User Login
// @Description  Authenticate user with email and password to get JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login Credentials"
// @Success      200 {object} response.DefaultResponse{data=LoginResponse}
// @Failure      400 {object} response.DefaultResponse
// @Failure      401 {object} response.DefaultResponse
// @Failure      403 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.service.Login(ctx, req)
	if err != nil {
		statusCode := http.StatusUnauthorized

		if errors.Is(err, ErrSystemFail) || errors.Is(err, ErrSessionFail) {
			statusCode = http.StatusInternalServerError
		} else if errors.Is(err, ErrAccountDisabled) {
			statusCode = http.StatusForbidden
		}

		response.Error(c, statusCode, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Login berhasil", res)
}
