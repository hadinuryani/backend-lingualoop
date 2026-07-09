package main

import (
	"log"
	"os"
	"path/filepath"

	"backend-lingualoop/config"
	"backend-lingualoop/database"
)

func main() {
	// Load configuration dari .env
	config.LoadConfig()

	// Koneksi ke database
	db := database.ConnectDB()
	defer db.Close()

	// Tentukan path folder migrations
	migrationsDir := filepath.Join("database", "migrations")

	// Cek apakah ada argumen custom path
	if len(os.Args) > 1 {
		migrationsDir = os.Args[1]
	}

	log.Println("================================================")
	log.Println("   LinguaLoop - Database Migration Runner")
	log.Println("================================================")
	log.Printf("📂 Folder migrasi: %s\n", migrationsDir)

	// Jalankan semua migrasi yang pending
	if err := database.RunMigrations(db, migrationsDir); err != nil {
		log.Fatalf("❌ Migration gagal: %v\n", err)
	}

	log.Println("================================================")
	log.Println("   Migration selesai!")
	log.Println("================================================")
}
