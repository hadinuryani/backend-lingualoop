package bootstrap

import (
	"database/sql"

	"backend-lingualoop/config"
	_ "backend-lingualoop/docs"
	"backend-lingualoop/internal/middleware"
	"backend-lingualoop/internal/modules/academic_year"
	"backend-lingualoop/internal/modules/auth"
	"backend-lingualoop/internal/modules/class"
	"backend-lingualoop/internal/modules/dashboard"
	"backend-lingualoop/internal/modules/major"
	"backend-lingualoop/internal/modules/student"
	"backend-lingualoop/internal/modules/subject"
	"backend-lingualoop/internal/modules/teacher"
	"backend-lingualoop/internal/modules/region"
	"backend-lingualoop/internal/modules/file"
	"backend-lingualoop/internal/modules/schedule"
	"backend-lingualoop/internal/modules/settings"
	"backend-lingualoop/internal/modules/teacher_portal"
	"backend-lingualoop/pkg/jwt"
	"backend-lingualoop/pkg/storage"

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

	// Static Route untuk file publik
	router.Static("/uploads/public", "./uploads/public")

	// 4. Setup API Grouping & Versioning
	api := router.Group("/api")
	v1 := api.Group("/v1")
	{
		// Registrasi Modul Fitur untuk V1
		cfg := config.GetConfig()
		jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.ExpirationHours)

		// Setup Storage
		baseURL := "http://localhost:" + cfg.App.Port // TODO: Get from env config for production
		store := storage.NewLocalStorage("./uploads", baseURL+"/uploads")

		// Register Public Route
		auth.RegisterRoute(v1, db, jwtManager)

		// Register Protected Admin Routes
		adminProtected := v1.Group("")
		adminProtected.Use(middleware.RequireAuth(jwtManager, db))
		adminProtected.Use(middleware.RequireRole("admin"))

		file.RegisterRoute(adminProtected, db, store)
		major.RegisterRoute(adminProtected, db, store)
		teacher.RegisterRoute(adminProtected, db)
		student.RegisterRoute(adminProtected, db)
		class.RegisterRoute(adminProtected, db)
		subject.RegisterRoute(adminProtected, db)
		schedule.RegisterRoute(adminProtected, db)
		academic_year.RegisterRoute(adminProtected, db)
		dashboard.RegisterRoute(adminProtected, db)
		region.RegisterRoute(adminProtected, db)
		settings.RegisterRoute(adminProtected, db, store)

		// Register Protected Teacher Routes
		teacherProtected := v1.Group("/teacher-portal")
		teacherProtected.Use(middleware.RequireAuth(jwtManager, db))
		teacherProtected.Use(middleware.RequireRole("teacher"))

		teacher_portal.RegisterRoute(teacherProtected, db)

	}

	// 5. Mount Swagger UI (Hanya jika bukan production)
	if !isProduction {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return router
}
