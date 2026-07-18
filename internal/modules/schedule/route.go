package schedule

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

func RegisterRoute(router *gin.RouterGroup, db *sql.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	schedules := router.Group("/schedules")
	{
		schedules.GET("/config", handler.GetConfig)
		schedules.PUT("/config", handler.SaveConfig)
		
		schedules.GET("", handler.GetAll)
		schedules.POST("", handler.Create)
		schedules.PUT("/:id", handler.Update)
		schedules.DELETE("/:id", handler.Delete)
		schedules.DELETE("/class/:class_id", handler.DeleteByClass)
	}
}
