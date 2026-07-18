package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"backend-lingualoop/pkg/security"
	"github.com/google/uuid"
)

func SeedUsers(ctx context.Context, db *sql.DB) error {
	slog.Info("Memeriksa dan membuat akun default (Admin, Teacher, Student)...")

	// Mulai Transaksi
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi seeder user: %w", err)
	}
	defer tx.Rollback()

	// Cleanup parsial (menangani inkonsistensi dari seeding sebelumnya)
	var exists bool
	_ = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM teachers LIMIT 1)").Scan(&exists)
	if !exists {
		_, _ = tx.ExecContext(ctx, "DELETE FROM users WHERE email = 'teacher@lingualoop.com'")
	}

	_ = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM students LIMIT 1)").Scan(&exists)
	if !exists {
		_, _ = tx.ExecContext(ctx, "DELETE FROM users WHERE email = 'student@lingualoop.com'")
	}

	query := `
		INSERT INTO users (id, email, username, password_hash, full_name, role, is_active)
		VALUES (?, ?, ?, ?, ?, ?, true)
	`

	_ = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE role = 'admin' LIMIT 1)").Scan(&exists)
	if !exists {
		adminID := uuid.New().String()
		email := "admin@lingualoop.com"
		username := "admin"
		fullName := "Administrator LinguaLoop"
		password := "admin123"
		hashedPassword, _ := security.HashPassword(password)

		_, err := tx.ExecContext(ctx, query, adminID, email, username, hashedPassword, fullName, "admin")
		if err != nil {
			return fmt.Errorf("gagal insert admin: %w", err)
		}
		slog.Info("Default Admin Created", "email", email, "password", password)
	}

	// 2. Konfigurasi Teacher
	_ = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE role = 'teacher' LIMIT 1)").Scan(&exists)
	if !exists {
		teacherUserID := uuid.New().String()
		tEmail := "teacher@lingualoop.com"
		tUsername := "teacher"
		tFullName := "Guru LinguaLoop"
		tPassword := "teacher123"
		tHashedPassword, _ := security.HashPassword(tPassword)

		_, err := tx.ExecContext(ctx, query, teacherUserID, tEmail, tUsername, tHashedPassword, tFullName, "teacher")
		if err != nil {
			return fmt.Errorf("gagal insert teacher ke users: %w", err)
		}

		teacherID := uuid.New().String()
		_, err = tx.ExecContext(ctx, "INSERT INTO teachers (id, user_id, nip, full_name, gender, status) VALUES (?, ?, ?, ?, ?, ?)",
			teacherID, teacherUserID, "T-10001", tFullName, "L", "ACTIVE")
		if err != nil {
			return fmt.Errorf("gagal insert teacher profile: %w", err)
		}

		slog.Info("Default Teacher Created", "email", tEmail, "password", tPassword)
	}

	// 3. Konfigurasi Student
	_ = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE role = 'student' LIMIT 1)").Scan(&exists)
	if !exists {
		studentUserID := uuid.New().String()
		sEmail := "student@lingualoop.com"
		sUsername := "student"
		sFullName := "Siswa LinguaLoop"
		sPassword := "student123"
		sHashedPassword, _ := security.HashPassword(sPassword)

		_, err := tx.ExecContext(ctx, query, studentUserID, sEmail, sUsername, sHashedPassword, sFullName, "student")
		if err != nil {
			return fmt.Errorf("gagal insert student ke users: %w", err)
		}

		studentID := uuid.New().String()
		_, err = tx.ExecContext(ctx, "INSERT INTO students (id, user_id, nis, full_name, gender, status) VALUES (?, ?, ?, ?, ?, ?)",
			studentID, studentUserID, "S-20001", sFullName, "P", "ACTIVE")
		if err != nil {
			return fmt.Errorf("gagal insert student profile: %w", err)
		}

		slog.Info("Default Student Created", "email", sEmail, "password", sPassword)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("gagal commit transaksi seeder users: %w", err)
	}

	slog.Info("Seeding selesai.")

	return nil
}
