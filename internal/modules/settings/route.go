package settings

import (
	"database/sql"

	"backend-lingualoop/pkg/storage"
	"github.com/gin-gonic/gin"
)

func RegisterRoute(router *gin.RouterGroup, db *sql.DB, store storage.Storage) {
	repo := NewRepository(db)
	service := NewService(repo, store)
	handler := NewHandler(service)

	s := router.Group("/settings")
	{
		s.GET("", handler.GetConfig)
		s.PUT("", handler.UpdateConfig)
	}
}
