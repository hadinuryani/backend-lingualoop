package main

import (
	"log"
	"path/filepath"

	"backend-lingualoop/config"
	"backend-lingualoop/database"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect ke database
	db := database.ConnectDB()
	defer db.Close()

	log.Println("   LinguaLoop - Database Seeder")

	// Jalankan seeder admin
	if err := database.SeedUsers(db); err != nil {
		log.Fatalf(" Gagal seeding user accounts: %v", err)
	}

	// Jalankan seeder wilayah Indonesia
	csvDir := filepath.Join("Wilayah-Indonesia-Beserta-Kode-Pos", "CSV")
	if err := database.SeedWilayah(db, csvDir); err != nil {
		log.Fatalf(" Gagal seeding wilayah Indonesia: %v", err)
	}

	log.Println("   Proses seeding selesai!")
}
