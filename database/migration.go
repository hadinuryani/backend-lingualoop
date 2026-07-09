package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Migration merepresentasikan satu file migrasi database
type Migration struct {
	Version  string
	Name     string
	FilePath string
}

// RunMigrations menjalankan semua file migrasi .sql yang belum dieksekusi
func RunMigrations(db *sql.DB, migrationsDir string) error {
	// 1. Buat tabel schema_migrations jika belum ada
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("gagal membuat tabel schema_migrations: %w", err)
	}

	// 2. Ambil daftar migrasi yang sudah dijalankan
	applied, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("gagal mengambil daftar migrasi: %w", err)
	}

	// 3. Scan folder migrations/ untuk file .sql
	migrations, err := scanMigrationFiles(migrationsDir)
	if err != nil {
		return fmt.Errorf("gagal scan file migrasi: %w", err)
	}

	if len(migrations) == 0 {
		log.Println("Tidak ada file migrasi ditemukan di:", migrationsDir)
		return nil
	}

	// 4. Filter hanya yang belum dijalankan
	pending := filterPending(migrations, applied)

	if len(pending) == 0 {
		log.Println("Database sudah up-to-date. Tidak ada migrasi baru.")
		return nil
	}

	// 5. Eksekusi migrasi yang pending secara berurutan
	log.Printf("Ditemukan %d migrasi yang perlu dijalankan...\n", len(pending))

	for _, m := range pending {
		log.Printf("▶ Menjalankan migrasi: %s (%s)\n", m.Version, m.Name)

		if err := executeMigration(db, m); err != nil {
			return fmt.Errorf("gagal menjalankan migrasi %s: %w", m.Version, err)
		}

		log.Printf(" Berhasil: %s\n", m.Version)
	}

	log.Printf("Semua %d migrasi berhasil dijalankan!\n", len(pending))
	return nil
}

// createMigrationsTable membuat tabel tracking migrasi jika belum ada
func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    VARCHAR(50) PRIMARY KEY,
			name       VARCHAR(255) NOT NULL,
			applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	_, err := db.Exec(query)
	return err
}

// getAppliedMigrations mengambil semua versi migrasi yang sudah dijalankan
func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	applied := make(map[string]bool)

	rows, err := db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return applied, err
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return applied, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// scanMigrationFiles membaca semua file .sql dari folder migrasi
func scanMigrationFiles(dir string) ([]Migration, error) {
	var migrations []Migration

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("tidak dapat membaca folder %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		// Format nama file: 001_create_users_table.sql
		filename := entry.Name()
		parts := strings.SplitN(strings.TrimSuffix(filename, ".sql"), "_", 2)

		version := parts[0]
		name := filename
		if len(parts) > 1 {
			name = parts[1]
		}

		migrations = append(migrations, Migration{
			Version:  version,
			Name:     name,
			FilePath: filepath.Join(dir, filename),
		})
	}

	// Sort berdasarkan version (nama file)
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// filterPending mengembalikan migrasi yang belum dijalankan
func filterPending(all []Migration, applied map[string]bool) []Migration {
	var pending []Migration
	for _, m := range all {
		if !applied[m.Version] {
			pending = append(pending, m)
		}
	}
	return pending
}

// executeMigration membaca file SQL dan menjalankannya dalam transaksi
func executeMigration(db *sql.DB, m Migration) error {
	// Baca isi file SQL
	content, err := os.ReadFile(m.FilePath)
	if err != nil {
		return fmt.Errorf("gagal membaca file %s: %w", m.FilePath, err)
	}

	sqlContent := string(content)
	if strings.TrimSpace(sqlContent) == "" {
		return fmt.Errorf("file migrasi kosong: %s", m.FilePath)
	}

	// Pisahkan per statement (berdasarkan delimiter ;)
	statements := splitStatements(sqlContent)

	// Jalankan dalam transaksi
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", err)
	}

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}

		if _, err := tx.Exec(stmt); err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal menjalankan statement: %w\nSQL: %s", err, truncate(stmt, 200))
		}
	}

	// Catat migrasi yang sudah dijalankan
	_, err = tx.Exec(
		"INSERT INTO schema_migrations (version, name, applied_at) VALUES (?, ?, ?)",
		m.Version, m.Name, time.Now(),
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal mencatat migrasi: %w", err)
	}

	return tx.Commit()
}

// splitStatements memecah string SQL menjadi statement-statement individual
func splitStatements(sql string) []string {
	var statements []string
	var current strings.Builder
	inString := false
	stringChar := byte(0)

	for i := 0; i < len(sql); i++ {
		ch := sql[i]

		// Handle string literals
		if !inString && (ch == '\'' || ch == '"') {
			inString = true
			stringChar = ch
			current.WriteByte(ch)
			continue
		}

		if inString && ch == stringChar {
			// Check for escaped quote
			if i+1 < len(sql) && sql[i+1] == stringChar {
				current.WriteByte(ch)
				current.WriteByte(ch)
				i++
				continue
			}
			inString = false
			current.WriteByte(ch)
			continue
		}

		// Handle single-line comments
		if !inString && ch == '-' && i+1 < len(sql) && sql[i+1] == '-' {
			for i < len(sql) && sql[i] != '\n' {
				i++
			}
			continue
		}

		// Handle statement delimiter
		if !inString && ch == ';' {
			stmt := strings.TrimSpace(current.String())
			if stmt != "" {
				statements = append(statements, stmt)
			}
			current.Reset()
			continue
		}

		current.WriteByte(ch)
	}

	// Sisa terakhir (tanpa ;)
	if stmt := strings.TrimSpace(current.String()); stmt != "" {
		statements = append(statements, stmt)
	}

	return statements
}

// truncate memotong string untuk logging
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
