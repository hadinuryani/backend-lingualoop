package middleware

import (
	"net/http"
	"strings"

	"backend-lingualoop/pkg/jwt"
	"backend-lingualoop/pkg/response"

	"github.com/gin-gonic/gin"
)

// RequireAuth mengecek keberadaan dan validitas JWT token di Header
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Akses ditolak: Token tidak ditemukan", nil)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "Akses ditolak: Format token tidak valid", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Akses ditolak: "+err.Error(), nil)
			c.Abort()
			return
		}

		// Simpan data user ke dalam context untuk dipakai oleh handler atau middleware berikutnya
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
