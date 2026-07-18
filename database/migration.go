package database

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"strconv"
)

type Migration struct {
	Version  int
	Name     string
	Checksum string
	FilePath string
}

var migrationFileRegex = regexp.MustCompile(`^([0-9]{3})_(.+)\.sql$`)

func RunMigrations(ctx context.Context, db *sql.DB, migrationsDir string) error {
	// 1. Ambil koneksi khusus untuk mengunci (GET_LOCK)
	conn, err := db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("gagal mendapatkan koneksi db: %w", err)
	}
	defer conn.Close()

	// Ambil lock
	var lockAcquired int
	err = conn.QueryRowContext(ctx, "SELECT GET_LOCK(CONCAT('lingualoop_migration_', DATABASE()), 10)").Scan(&lockAcquired)
	if err != nil {
		return fmt.Errorf("gagal mencoba acquire lock: %w", err)
	}
	if lockAcquired != 1 {
		return fmt.Errorf("migration lock sedang digunakan oleh proses lain, harap tunggu")
	}
	// Pastikan lock dilepas saat fungsi selesai
	defer func() {
		_, _ = conn.ExecContext(ctx, "SELECT RELEASE_LOCK(CONCAT('lingualoop_migration_', DATABASE()))")
	}()

	slog.Info("Migration lock acquired")

	// 2. Buat/Update tabel schema_migrations
	if err := createMigrationsTable(conn, ctx); err != nil {
		return fmt.Errorf("gagal membuat tabel schema_migrations: %w", err)
	}

	// 3. Ambil daftar migrasi yang sudah dijalankan
	applied, err := getAppliedMigrations(conn, ctx)
	if err != nil {
		return fmt.Errorf("gagal mengambil daftar migrasi: %w", err)
	}

	// 4. Scan folder migrations/
	migrations, err := scanMigrationFiles(migrationsDir)
	if err != nil {
		return fmt.Errorf("gagal scan file migrasi: %w", err)
	}

	if len(migrations) == 0 {
		slog.Info("Tidak ada file migrasi ditemukan di: " + migrationsDir)
		return nil
	}

	// 5. Verifikasi Integritas (Duplicate, Checksum, Missing)
	if err := verifyMigrations(migrations, applied); err != nil {
		return fmt.Errorf("verifikasi migrasi gagal: %w", err)
	}

	// 6. Filter hanya yang belum dijalankan
	pending := filterPending(migrations, applied)

	if len(pending) == 0 {
		slog.Info("Database sudah up-to-date. Tidak ada migrasi baru.")
		return nil
	}

	// 7. Eksekusi migrasi yang pending secara berurutan
	slog.Info(fmt.Sprintf("Ditemukan %d migrasi yang perlu dijalankan...", len(pending)))

	for _, m := range pending {
		slog.Info(fmt.Sprintf("Menjalankan migrasi: %03d (%s)", m.Version, m.Name))

		if err := executeMigration(conn, ctx, m); err != nil {
			return fmt.Errorf("gagal menjalankan migrasi %03d: %w", m.Version, err)
		}

		slog.Info(fmt.Sprintf("Berhasil: %03d", m.Version))
	}

	slog.Info(fmt.Sprintf("Semua %d migrasi berhasil dijalankan!", len(pending)))
	return nil
}

func createMigrationsTable(conn *sql.Conn, ctx context.Context) error {
	// Pengecekan tipe data kolom version untuk memutuskan perlunya Drop
	var columnType string
	err := conn.QueryRowContext(ctx, "SELECT DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'schema_migrations' AND COLUMN_NAME = 'version' AND TABLE_SCHEMA = DATABASE()").Scan(&columnType)
	
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	
	// Jika tabel sudah ada dan kolom version bukan int (yaitu varchar peninggalan lama), lakukan ALTER TABLE
	if err == nil && !strings.Contains(strings.ToLower(columnType), "int") {
		slog.Warn("Tabel schema_migrations lama terdeteksi (tipe varchar). Akan dilakukan ALTER TABLE untuk mengupgrade skema migration runner.")
		
		alterQuery := `
			ALTER TABLE schema_migrations 
			DROP PRIMARY KEY,
			MODIFY COLUMN version INT,
			ADD PRIMARY KEY (version),
			ADD COLUMN checksum CHAR(64) NOT NULL DEFAULT '' AFTER name;
		`
		if _, err := conn.ExecContext(ctx, alterQuery); err != nil {
			return fmt.Errorf("gagal melakukan ALTER TABLE pada schema_migrations: %w", err)
		}
	}

	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    INT PRIMARY KEY,
			name       VARCHAR(255) NOT NULL,
			checksum   CHAR(64) NOT NULL,
			applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	_, err = conn.ExecContext(ctx, query)
	return err
}

func getAppliedMigrations(conn *sql.Conn, ctx context.Context) (map[int]Migration, error) {
	applied := make(map[int]Migration)

	rows, err := conn.QueryContext(ctx, "SELECT version, name, checksum FROM schema_migrations ORDER BY version")
	if err != nil {
		return applied, err
	}
	defer rows.Close()

	for rows.Next() {
		var m Migration
		if err := rows.Scan(&m.Version, &m.Name, &m.Checksum); err != nil {
			return applied, err
		}
		applied[m.Version] = m
	}

	return applied, rows.Err()
}

func scanMigrationFiles(dir string) ([]Migration, error) {
	var migrations []Migration

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("tidak dapat membaca folder %s: %w", dir, err)
	}

	seenVersions := make(map[int]string)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		filename := entry.Name()
		matches := migrationFileRegex.FindStringSubmatch(filename)
		if len(matches) != 3 {
			return nil, fmt.Errorf("format nama file migrasi tidak valid: %s (Harus sesuai regex ^([0-9]{3})_(.+)\\.sql$)", filename)
		}

		versionStr := matches[1]
		name := matches[2]
		
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			return nil, fmt.Errorf("gagal parsing versi migrasi dari file %s", filename)
		}

		if existingFile, exists := seenVersions[version]; exists {
			return nil, fmt.Errorf("DUPLICATE MIGRATION VERSION TERDETEKSI: Versi %03d dipakai oleh '%s' dan '%s'", version, existingFile, filename)
		}
		seenVersions[version] = filename

		filePath := filepath.Join(dir, filename)
		checksum, err := calculateChecksum(filePath)
		if err != nil {
			return nil, fmt.Errorf("gagal menghitung checksum file %s: %w", filename, err)
		}

		migrations = append(migrations, Migration{
			Version:  version,
			Name:     name,
			Checksum: checksum,
			FilePath: filePath,
		})
	}

	// Sort berdasarkan integer version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func calculateChecksum(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:]), nil
}

func verifyMigrations(migrations []Migration, applied map[int]Migration) error {
	var maxAppliedVersion int
	for v := range applied {
		if v > maxAppliedVersion {
			maxAppliedVersion = v
		}
	}

	// Cek out-of-order dari sisi file yang akan di-apply
	for _, m := range migrations {
		if appliedM, exists := applied[m.Version]; exists {
			// Sudah di-apply, cek apakah file dimodifikasi
			// Hanya bandingkan checksum jika appliedM.Checksum tidak kosong (bisa jadi kosong akibat dari ALTER TABLE versi lama)
			if appliedM.Checksum != "" && appliedM.Checksum != m.Checksum {
				return fmt.Errorf("MODIFIED MIGRATION TERDETEKSI: File '%s' (Versi %d) telah diubah setelah dijalankan. Checksum DB: %s, File: %s", m.FilePath, m.Version, appliedM.Checksum, m.Checksum)
			}
		} else {
			// Belum di-apply, cek apakah ini out-of-order
			if m.Version < maxAppliedVersion {
				return fmt.Errorf("MISSING MIGRATION TERDETEKSI: Versi %d ('%s') terlewat dan ada versi lebih baru yang sudah berjalan di database", m.Version, m.FilePath)
			}
		}
	}

	// Cek missing file: ada di database tapi file fisiknya hilang
	fileMap := make(map[int]Migration)
	for _, m := range migrations {
		fileMap[m.Version] = m
	}
	for v, appliedM := range applied {
		if _, exists := fileMap[v]; !exists {
			return fmt.Errorf("MISSING FILE TERDETEKSI: Migrasi versi %03d ('%s') tercatat di database tapi file fisiknya tidak ditemukan", v, appliedM.Name)
		}
	}

	return nil
}

func filterPending(all []Migration, applied map[int]Migration) []Migration {
	var pending []Migration
	for _, m := range all {
		if _, exists := applied[m.Version]; !exists {
			pending = append(pending, m)
		}
	}
	return pending
}

func executeMigration(conn *sql.Conn, ctx context.Context, m Migration) error {
	content, err := os.ReadFile(m.FilePath)
	if err != nil {
		return fmt.Errorf("gagal membaca file %s: %w", m.FilePath, err)
	}

	sqlContent := string(content)
	if strings.TrimSpace(sqlContent) == "" {
		return fmt.Errorf("file migrasi kosong: %s", m.FilePath)
	}

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi migrasi: %w", err)
	}

	// Eksekusi keseluruhan isi file
	// MySQL Driver akan mengeksekusi multi-statement karena kita menambah multiStatements=true di DSN
	_, err = tx.ExecContext(ctx, sqlContent)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menjalankan query migrasi: %w", err)
	}

	// Catat migrasi yang sudah dijalankan
	_, err = tx.ExecContext(ctx, 
		"INSERT INTO schema_migrations (version, name, checksum, applied_at) VALUES (?, ?, ?, ?)",
		m.Version, m.Name, m.Checksum, time.Now(),
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal mencatat migrasi ke schema_migrations: %w", err)
	}

	return tx.Commit()
}
