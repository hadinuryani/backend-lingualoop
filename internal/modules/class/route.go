package class

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

func RegisterRoute(router *gin.RouterGroup, db *sql.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	classes := router.Group("/classes")
	{
		classes.GET("", handler.GetAll)
		classes.GET("/:id", handler.GetByID)
		classes.POST("", handler.Create)
		classes.POST("/batch", handler.CreateBatch)
		classes.PUT("/:id", handler.Update)
		classes.DELETE("/:id", handler.Delete)
	}
}
