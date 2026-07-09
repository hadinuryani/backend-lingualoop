package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler menangani HTTP request dan response untuk modul auth
type Handler struct {
	service Service
}

// NewHandler membuat instance baru dari auth handler
func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Login menangani request endpoint POST /login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	// 1. Parse dan Validasi JSON Input (Tugas Handler)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format email atau password tidak valid"})
		return
	}

	// 2. Lempar ke Service untuk diproses (Business Logic)
	res, err := h.service.Login(req)

	// 3. Tangani hasil dari service (Tugas Handler)
	if err != nil {
		// Menggunakan StatusUnauthorized untuk error credential (bisa disesuaikan di service if needed)
		statusCode := http.StatusUnauthorized
		if err.Error() == "terjadi kesalahan pada sistem, silakan coba lagi" {
			statusCode = http.StatusInternalServerError
		} else if err.Error() == "akun Anda dinonaktifkan, silakan hubungi admin" {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	// 4. Kembalikan JSON sukses (Tugas Handler)
	c.JSON(http.StatusOK, res)
}
