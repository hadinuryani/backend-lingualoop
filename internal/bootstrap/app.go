package bootstrap

import (
	"database/sql"

	"backend-lingualoop/config"
	_ "backend-lingualoop/docs"
	"backend-lingualoop/internal/middleware"
	"backend-lingualoop/internal/modules/academic_year"
	"backend-lingualoop/internal/modules/auth"
	"backend-lingualoop/internal/modules/class"
	"backend-lingualoop/internal/modules/major"
	"backend-lingualoop/internal/modules/student"
	"backend-lingualoop/internal/modules/subject"
	"backend-lingualoop/internal/modules/teacher"
	"backend-lingualoop/pkg/jwt"

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
		cfg := config.GetConfig()
		jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.ExpirationHours)

		// Register Public Route
		auth.RegisterRoute(v1, db, jwtManager)

		// Register Protected Admin Routes
		adminProtected := v1.Group("")
		adminProtected.Use(middleware.RequireAuth(jwtManager, db))
		adminProtected.Use(middleware.RequireRole("admin"))

		major.RegisterRoute(adminProtected, db)
		teacher.RegisterRoute(adminProtected, db)
		student.RegisterRoute(adminProtected, db)
		class.RegisterRoute(adminProtected, db)
		subject.RegisterRoute(adminProtected, db)
		academic_year.RegisterRoute(adminProtected, db)

	}

	// 5. Mount Swagger UI (Hanya jika bukan production)
	if !isProduction {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return router
}
