package main

import (
	"log"

	"backend-lingualoop/config"
	"backend-lingualoop/database"
	"backend-lingualoop/internal/bootstrap"

	"github.com/gin-gonic/gin"
)

// @title           LinguaLoop API
// @version         1.0
// @description     API Documentation for LinguaLoop E-Learning Platform.
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 1. Load Konfigurasi (.env)
	cfg := config.LoadConfig()

	// 2. Konek ke Database MySQL
	db := database.ConnectDB()
	defer db.Close()

	// 3. Inisialisasi Aplikasi (Router & Semua Modul via Bootstrap)
	isProduction := cfg.App.Env == "production"
	app := bootstrap.SetupApp(db, isProduction)

	// 4. Endpoint Health Check
	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"app":    cfg.App.Name,
		})
	})

	// 5. Jalankan Server
	log.Printf("Starting %s on port %s\n", cfg.App.Name, cfg.App.Port)
	log.Printf("Environment: %s\n", cfg.App.Env)
	log.Printf("Database: %s@%s:%s\n", cfg.Database.User, cfg.Database.Host, cfg.Database.Port)

	if err := app.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
