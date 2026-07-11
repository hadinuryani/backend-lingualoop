package middleware

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strings"

	"backend-lingualoop/pkg/jwt"
	"backend-lingualoop/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	CtxUserID   = "user_id"
	CtxEmail    = "email"
	CtxUsername = "username"
	CtxRole     = "role"
)

// RequireAuth mengecek keberadaan dan validitas JWT token di Header
func RequireAuth(jwtManager jwt.Manager, db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Token tidak ditemukan", nil)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(c, http.StatusUnauthorized, "Format token tidak valid", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			slog.Error("JWT Validation failed", "error", err)
			response.Error(c, http.StatusUnauthorized, "Token tidak valid", nil)
			c.Abort()
			return
		}

		// Cek apakah user masih aktif di database
		var isActive bool
		err = db.QueryRowContext(c.Request.Context(), "SELECT is_active FROM users WHERE id = ? AND deleted_at IS NULL", claims.UserID).Scan(&isActive)
		if err != nil {
			if err == sql.ErrNoRows {
				slog.Error("Auth failed: user not found", "user_id", claims.UserID)
				response.Error(c, http.StatusUnauthorized, "Akun tidak ditemukan", nil)
				c.Abort()
				return
			}
			slog.Error("Database query failed during auth", "error", err, "user_id", claims.UserID)
			response.Error(c, http.StatusInternalServerError, "Terjadi kesalahan sistem", nil)
			c.Abort()
			return
		}

		if !isActive {
			slog.Error("Auth failed: user inactive", "user_id", claims.UserID)
			response.Error(c, http.StatusForbidden, "Akun Anda telah dinonaktifkan", nil)
			c.Abort()
			return
		}

		// Simpan data user ke dalam context untuk dipakai oleh handler atau middleware berikutnya
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxEmail, claims.Email)
		c.Set(CtxUsername, claims.Username)
		c.Set(CtxRole, claims.Role)

		c.Next()
	}
}
