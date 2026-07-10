package student

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

func RegisterRoute(router *gin.RouterGroup, db *sql.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	students := router.Group("/students")
	{
		students.GET("", handler.GetAll)
		students.GET("/:id", handler.GetByID)
		students.POST("", handler.Create)
		students.PUT("/:id", handler.Update)
		students.PATCH("/:id/status", handler.UpdateStatus)
		students.DELETE("/:id", handler.Delete)
	}
}
