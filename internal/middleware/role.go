package middleware

import (
	"net/http"

	"backend-lingualoop/pkg/response"

	"github.com/gin-gonic/gin"
)

// RequireRole membatasi akses hanya untuk role tertentu
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "Akses ditolak: Role tidak teridentifikasi", nil)
			c.Abort()
			return
		}

		userRole, ok := roleVal.(string)
		if !ok {
			response.Error(c, http.StatusInternalServerError, "Kesalahan server internal: Tipe role tidak valid", nil)
			c.Abort()
			return
		}

		// Cek apakah userRole ada di dalam daftar allowedRoles
		isAllowed := false
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			response.Error(c, http.StatusForbidden, "Akses ditolak: Anda tidak memiliki wewenang untuk endpoint ini", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
