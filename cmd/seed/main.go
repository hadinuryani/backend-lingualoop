package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"backend-lingualoop/config"
	"backend-lingualoop/database"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect ke database
	db, err := database.ConnectDB()
	if err != nil {
		slog.Error("Database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("LinguaLoop - Database Seeder")

	// Jalankan seeder admin
	if err := database.SeedUsers(context.Background(), db); err != nil {
		slog.Error("Gagal seeding user accounts", "error", err)
		os.Exit(1)
	}

	// Jalankan seeder wilayah Indonesia
	csvDir := filepath.Join("Wilayah-Indonesia-Beserta-Kode-Pos", "CSV")
	if err := database.SeedWilayah(context.Background(), db, csvDir); err != nil {
		slog.Error("Gagal seeding wilayah Indonesia", "error", err)
		os.Exit(1)
	}

	slog.Info("Proses seeding selesai!")
}
