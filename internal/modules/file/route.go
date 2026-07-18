package file

import (
	"database/sql"
	"backend-lingualoop/pkg/storage"
	"github.com/gin-gonic/gin"
)

func RegisterRoute(router *gin.RouterGroup, db *sql.DB, store storage.Storage) {
	repo := NewRepository(db)
	service := NewService(repo, store)
	handler := NewHandler(service)

	fileRoute := router.Group("/files")
	{
		fileRoute.POST("/upload-image", handler.UploadImage)
	}
}
