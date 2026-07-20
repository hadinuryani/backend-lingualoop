package teacher_portal

import (
	"context"
	"net/http"
	"time"

	"backend-lingualoop/internal/middleware"
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

// GetMyClasses mengambil daftar kelas yang diajar oleh guru yang sedang login
// @Summary      Get teacher's classes
// @Description  Get a list of classes assigned to the currently authenticated teacher
// @Tags         Teacher Portal
// @Produce      json
// @Success      200 {object} response.DefaultResponse{data=[]TeacherClassResponse}
// @Failure      401 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /teacher-portal/classes [get]
// @Security     BearerAuth
func (h *Handler) GetMyClasses(c *gin.Context) {
	userID := c.GetString(middleware.CtxUserID)
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, "User ID tidak ditemukan di context", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	classes, err := h.service.GetMyClasses(ctx, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Handle empty array to return [] instead of null in JSON
	if classes == nil {
		classes = []TeacherClassResponse{}
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil kelas Anda", classes)
}

// GetMySchedules mengambil jadwal mengajar (hari ini & mingguan) untuk guru yang sedang login
// @Summary      Get teacher's schedules
// @Description  Get today's and weekly schedule for the authenticated teacher
// @Tags         Teacher Portal
// @Produce      json
// @Success      200 {object} response.DefaultResponse
// @Failure      401 {object} response.DefaultResponse
// @Failure      500 {object} response.DefaultResponse
// @Router       /teacher-portal/schedules [get]
// @Security     BearerAuth
func (h *Handler) GetMySchedules(c *gin.Context) {
	userID := c.GetString(middleware.CtxUserID)
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, "User ID tidak ditemukan di context", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schedules, err := h.service.GetMySchedules(ctx, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil jadwal Anda", schedules)
}
