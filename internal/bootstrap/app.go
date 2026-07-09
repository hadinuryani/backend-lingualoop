package bootstrap

import (
	"database/sql"

	"backend-lingualoop/internal/middleware"
	"backend-lingualoop/internal/modules/auth"

	"github.com/gin-gonic/gin"
)

// SetupApp menginisialisasi router Gin, mendaftarkan global middleware,
// dan merangkai (wiring) seluruh rute dari setiap modul fitur.
func SetupApp(db *sql.DB, isProduction bool) *gin.Engine {
	// 1. Set mode framework
	if isProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	// 2. Inisialisasi router
	router := gin.Default()

	// 3. Pasang Global Middleware (seperti CORS)
	router.Use(middleware.CORS())

	// 4. Setup API Grouping & Versioning
	api := router.Group("/api")
	v1 := api.Group("/v1")
	{
		// Registrasi Modul Fitur untuk V1
		auth.RegisterRoute(v1, db)

		// Modul lain di masa depan:
		// assignment.RegisterRoute(v1, db)
		// grade.RegisterRoute(v1, db)
		// material.RegisterRoute(v1, db)
	}

	return router
}
