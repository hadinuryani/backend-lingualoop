package main

import (
	"fmt"
	"log/slog"
	"os"

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
const banner = "\n" +
"     ____             __                  __   ___                         __                \n" +
"   / __ )____ ______/ /_____  ____  ____/ /  / (_)___  ____ ___  ______ _/ /___  ____  ____ \n" +
"  / __  / __ `/ ___/ //_/ _ \\/ __ \\/ __  /  / / / __ \\/ __ `/ / / / __ `/ / __ \\/ __ \\/ __ \\\n" +
" / /_/ / /_/ / /__/ ,< /  __/ / / / /_/ /  / / / / / / /_/ / /_/ / /_/ / / /_/ / /_/ / /_/ /\n" +
"/_____/\\__,_/\\___/_/|_|\\___/_/ /_/\\__,_/  /_/_/_/ /_/\\__, /\\__,_/\\__,_/_/\\____/\\____/ .___/ \n" +
"                                                    /____/                         /_/      \n"

func main() {
	fmt.Println(banner)

	// 1. Load Konfigurasi (.env)
	cfg := config.LoadConfig()

	// 2. Konek ke Database MySQL
	db, err := database.ConnectDB()
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	
	// Set instance global agar GetDB() dari package lain bisa mendapatkan instance yang sama
	database.DB = db

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
	slog.Info(fmt.Sprintf("Starting %s on port %s", cfg.App.Name, cfg.App.Port))
	slog.Info(fmt.Sprintf("Environment: %s", cfg.App.Env))
	slog.Info(fmt.Sprintf("Database: %s@%s:%s", cfg.Database.User, cfg.Database.Host, cfg.Database.Port))

	if err := app.Run(":" + cfg.App.Port); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
