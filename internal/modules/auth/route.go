package auth

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoute(router *gin.RouterGroup, db *sql.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	// 2. Daftarkan routes
	authRoute := router.Group("/auth")
	{
		authRoute.POST("/login", handler.Login)
		// authRoute.GET("/me", middleware.Auth(), handler.Me) -> contoh nanti
	}
}
