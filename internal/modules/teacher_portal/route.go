package teacher_portal

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

func RegisterRoute(router *gin.RouterGroup, db *sql.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	router.GET("/classes", handler.GetMyClasses)
	router.GET("/schedules", handler.GetMySchedules)
}
