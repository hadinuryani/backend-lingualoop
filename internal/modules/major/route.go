package major

import (
	"database/sql"

	"backend-lingualoop/pkg/storage"
	"github.com/gin-gonic/gin"
)

func RegisterRoute(router *gin.RouterGroup, db *sql.DB, store storage.Storage) {
	repo := NewRepository(db)
	service := NewService(repo, store)
	handler := NewHandler(service)

	majorRoute := router.Group("/majors")
	{
		majorRoute.GET("", handler.GetAll)
		majorRoute.POST("", handler.Create)
		majorRoute.PUT("/:id", handler.Update)
		majorRoute.DELETE("/:id", handler.Delete)
	}
}
