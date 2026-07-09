package database

import (
	"database/sql"
	"fmt"
	"log"

	"backend-lingualoop/pkg/security"
	"github.com/google/uuid"
)

func SeedAdmin(db *sql.DB) error {
	var count int
	err := db.QueryRow("SELECT COUNT(id) FROM users WHERE role = 'admin'").Scan(&count)
	if err != nil {
		return fmt.Errorf("gagal mengecek admin: %w", err)
	}

	if count > 0 {
		log.Println("Akun admin sudah ada di database, melewati proses seeding.")
		return nil
	}

	log.Println("Membuat akun admin default...")

	// Konfigurasi admin default
	adminID := uuid.New().String()
	email := "admin@lingualoop.com"
	username := "admin"
	fullName := "Administrator LinguaLoop"
	password := "admin123"

	// Hash password
	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		return fmt.Errorf("gagal hash password admin: %w", err)
	}

	// Insert ke database
	query := `
		INSERT INTO users (id, email, username, password_hash, full_name, role, is_active)
		VALUES (?, ?, ?, ?, ?, 'admin', true)
	`
	_, err = db.Exec(query, adminID, email, username, hashedPassword, fullName)
	if err != nil {
		return fmt.Errorf("gagal insert admin: %w", err)
	}

	log.Println(" Akun admin berhasil dibuat!")
	log.Printf("Email: %s\n", email)
	log.Printf("Password: %s\n", password)

	return nil
}