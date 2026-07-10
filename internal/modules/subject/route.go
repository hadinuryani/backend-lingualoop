package subject

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

func RegisterRoute(router *gin.RouterGroup, db *sql.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	subjects := router.Group("/subjects")
	{
		subjects.GET("", handler.GetAll)
		subjects.GET("/:id", handler.GetByID)
		subjects.POST("", handler.Create)
		subjects.PUT("/:id", handler.Update)
		subjects.DELETE("/:id", handler.Delete)
	}
}
