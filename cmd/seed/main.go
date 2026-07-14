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

	log.Println("   LinguaLoop - Database Seeder")

	// Jalankan seeder admin
	if err := database.SeedUsers(db); err != nil {
		log.Fatalf(" Gagal seeding user accounts: %v", err)
	}

	log.Println("   Proses seeding selesai!")
}
