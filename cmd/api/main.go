package main

import (
	"log"

	"backend-lingualoop/config"
	"backend-lingualoop/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration dari .env
	cfg := config.LoadConfig()

	// Inisialisasi koneksi ke database MySQL
	database.ConnectDB()

	// Set Gin mode berdasarkan environment
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"app":    cfg.App.Name,
		})
	})

	// Print configuration info
	log.Printf("Starting %s on port %s\n", cfg.App.Name, cfg.App.Port)
	log.Printf("Environment: %s\n", cfg.App.Env)
	log.Printf("Database: %s@%s:%s\n", cfg.Database.User, cfg.Database.Host, cfg.Database.Port)

	// Run server
	if err := router.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
