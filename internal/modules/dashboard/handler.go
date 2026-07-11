package dashboard

import (
	"context"
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

// GetDashboardStats godoc
// @Summary      Get Admin Dashboard Stats
// @Description  Get aggregated statistics for admin dashboard including demograhics and recent activities.
// @Tags         Dashboard
// @Produce      json
// @Success      200 {object} response.DefaultResponse{data=DashboardResponse}
// @Failure      500 {object} response.DefaultResponse
// @Router       /dashboard/stats [get]
// @Security     BearerAuth
func (h *Handler) GetDashboardStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	data, err := h.service.GetDashboardData(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil data dashboard", nil)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil data dashboard", data)
}
