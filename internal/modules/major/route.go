package major

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoute(router *gin.RouterGroup, db *sql.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	majorRoute := router.Group("/majors")
	{
		majorRoute.GET("", handler.GetAll)
		majorRoute.POST("", handler.Create)
		majorRoute.PUT("/:id", handler.Update)
		majorRoute.DELETE("/:id", handler.Delete)
	}
}
