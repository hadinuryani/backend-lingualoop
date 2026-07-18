package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"backend-lingualoop/config"
	"backend-lingualoop/database"
)

func main() {
	config.LoadConfig()

	db, err := database.ConnectDB()
	if err != nil {
		slog.Error("Database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Tentukan path folder migrations
	migrationsDir := filepath.Join("database", "migrations")

	// Cek apakah ada argumen custom path
	if len(os.Args) > 1 {
		migrationsDir = os.Args[1]
	}

	slog.Info("LinguaLoop - Database Migration Runner")
	slog.Info(fmt.Sprintf("Folder migrasi: %s", migrationsDir))

	// Jalankan semua migrasi yang pending
	if err := database.RunMigrations(context.Background(), db, migrationsDir); err != nil {
		slog.Error("Migration gagal", "error", err)
		os.Exit(1)
	}
	slog.Info("Migration selesai!")
}
