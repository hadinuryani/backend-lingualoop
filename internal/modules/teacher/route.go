package teacher

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

func RegisterRoute(router *gin.RouterGroup, db *sql.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	teachers := router.Group("/teachers")
	{
		teachers.GET("", handler.GetAll)
		teachers.GET("/:id", handler.GetByID)
		teachers.POST("", handler.Create)
		teachers.PUT("/:id", handler.Update)
		teachers.PATCH("/:id/status", handler.ToggleStatus)
		teachers.DELETE("/:id", handler.Delete)
	}
}
