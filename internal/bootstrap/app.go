package bootstrap

import (
	"database/sql"

	_ "backend-lingualoop/docs" 
	"backend-lingualoop/internal/middleware"
	"backend-lingualoop/internal/modules/auth"
	"backend-lingualoop/internal/modules/class"
	"backend-lingualoop/internal/modules/major"
	"backend-lingualoop/internal/modules/student"
	"backend-lingualoop/internal/modules/subject"
	"backend-lingualoop/internal/modules/teacher"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

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
		major.RegisterRoute(v1, db)
		teacher.RegisterRoute(v1, db)
		student.RegisterRoute(v1, db)
		class.RegisterRoute(v1, db)
		subject.RegisterRoute(v1, db)


	}

	// 5. Mount Swagger UI (Hanya jika bukan production)
	if !isProduction {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return router
}
