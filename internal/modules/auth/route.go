package auth

import (
	"backend-lingualoop/pkg/jwt"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoute(router *gin.RouterGroup, db *sql.DB, jwtManager jwt.Manager) {
	repo := NewRepository(db)
	service := NewService(repo, jwtManager)
	handler := NewHandler(service)

	// 2. Daftarkan routes
	authRoute := router.Group("/auth")
	{
		authRoute.POST("/login", handler.Login)
		// authRoute.GET("/me", middleware.Auth(), handler.Me) -> contoh nanti
	}
}
