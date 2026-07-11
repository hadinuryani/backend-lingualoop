package dashboard

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoute(router *gin.RouterGroup, db *sql.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	dashboardRoute := router.Group("/dashboard")
	{
		dashboardRoute.GET("/stats", handler.GetDashboardStats)
	}
}
