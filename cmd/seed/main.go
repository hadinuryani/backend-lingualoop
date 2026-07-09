package main

import (
	"log"

	"backend-lingualoop/config"
	"backend-lingualoop/database"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect ke database
	db := database.ConnectDB()
	defer db.Close()

	log.Println("================================================")
	log.Println("   LinguaLoop - Database Seeder")
	log.Println("================================================")

	// Jalankan seeder admin
	if err := database.SeedAdmin(db); err != nil {
		log.Fatalf("❌ Seeder gagal: %v\n", err)
	}

	log.Println("================================================")
	log.Println("   Proses seeding selesai!")
	log.Println("================================================")
}
