package database

import (
	"database/sql"
	"fmt"
	"log"

	"backend-lingualoop/pkg/security"
	"github.com/google/uuid"
)

func SeedUsers(db *sql.DB) error {
	log.Println("Memeriksa dan membuat akun default (Admin, Teacher, Student)...")

	// Cleanup parsial (menangani inkonsistensi dari seeding sebelumnya)
	var teacherProfileCount int
	db.QueryRow("SELECT COUNT(id) FROM teachers").Scan(&teacherProfileCount)
	if teacherProfileCount == 0 {
		db.Exec("DELETE FROM users WHERE email = 'teacher@lingualoop.com'")
	}

	var studentProfileCount int
	db.QueryRow("SELECT COUNT(id) FROM students").Scan(&studentProfileCount)
	if studentProfileCount == 0 {
		db.Exec("DELETE FROM users WHERE email = 'student@lingualoop.com'")
	}

	query := `
		INSERT INTO users (id, email, username, password_hash, full_name, role, is_active)
		VALUES (?, ?, ?, ?, ?, ?, true)
	`

	var adminCount int
	db.QueryRow("SELECT COUNT(id) FROM users WHERE role = 'admin'").Scan(&adminCount)
	if adminCount == 0 {
		adminID := uuid.New().String()
		email := "admin@lingualoop.com"
		username := "admin"
		fullName := "Administrator LinguaLoop"
		password := "admin123"
		hashedPassword, _ := security.HashPassword(password)

		_, err := db.Exec(query, adminID, email, username, hashedPassword, fullName, "admin")
		if err != nil {
			return fmt.Errorf("gagal insert admin: %w", err)
		}
		log.Printf("Admin   : %s / %s\n", email, password)
	}

	// 2. Konfigurasi Teacher
	var teacherCount int
	db.QueryRow("SELECT COUNT(id) FROM users WHERE role = 'teacher'").Scan(&teacherCount)
	if teacherCount == 0 {
		teacherUserID := uuid.New().String()
		tEmail := "teacher@lingualoop.com"
		tUsername := "teacher"
		tFullName := "Guru LinguaLoop"
		tPassword := "teacher123"
		tHashedPassword, _ := security.HashPassword(tPassword)

		_, err := db.Exec(query, teacherUserID, tEmail, tUsername, tHashedPassword, tFullName, "teacher")
		if err != nil {
			return fmt.Errorf("gagal insert teacher ke users: %w", err)
		}

		teacherID := uuid.New().String()
		_, err = db.Exec("INSERT INTO teachers (id, user_id, nip, full_name, gender, status) VALUES (?, ?, ?, ?, ?, ?)",
			teacherID, teacherUserID, "T-10001", tFullName, "L", "ACTIVE")
		if err != nil {
			return fmt.Errorf("gagal insert teacher profile: %w", err)
		}

		log.Printf("Teacher : %s / %s\n", tEmail, tPassword)
	}

	// 3. Konfigurasi Student
	var studentCount int
	db.QueryRow("SELECT COUNT(id) FROM users WHERE role = 'student'").Scan(&studentCount)
	if studentCount == 0 {
		studentUserID := uuid.New().String()
		sEmail := "student@lingualoop.com"
		sUsername := "student"
		sFullName := "Siswa LinguaLoop"
		sPassword := "student123"
		sHashedPassword, _ := security.HashPassword(sPassword)

		_, err := db.Exec(query, studentUserID, sEmail, sUsername, sHashedPassword, sFullName, "student")
		if err != nil {
			return fmt.Errorf("gagal insert student ke users: %w", err)
		}

		studentID := uuid.New().String()
		_, err = db.Exec("INSERT INTO students (id, user_id, nis, full_name, gender, status) VALUES (?, ?, ?, ?, ?, ?)",
			studentID, studentUserID, "S-20001", sFullName, "P", "ACTIVE")
		if err != nil {
			return fmt.Errorf("gagal insert student profile: %w", err)
		}

		log.Printf("Student : %s / %s\n", sEmail, sPassword)
	}

	log.Println("Seeding selesai.")

	return nil
}
