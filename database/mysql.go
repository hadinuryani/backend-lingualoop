package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"backend-lingualoop/config"

	_ "github.com/go-sql-driver/mysql"
)

// DB adalah instance koneksi database global
var (
	DB   *sql.DB
	once sync.Once
)

// ConnectDB membuka koneksi ke database dan mengembalikan db instance atau error
func ConnectDB() (*sql.DB, error) {
	cfg := config.GetConfig()

	// DSN eksplisit dengan konfigurasi lengkap (sebagaimana best-practice)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&loc=Local&multiStatements=true&timeout=5s&readTimeout=30s&writeTimeout=30s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi sql: %w", err)
	}

	// Konfigurasi Connection Pool
	db.SetMaxOpenConns(cfg.Database.MaxConnections)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConnections)
	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(15 * time.Minute)

	// Menggunakan PingContext dengan timeout agar aplikasi tidak hang jika DB mati
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("gagal terhubung ke database (Ping timeout): %w", err)
	}

	slog.Info("Koneksi ke database MySQL berhasil dibentuk")
	return db, nil
}

// GetDB mengembalikan instance koneksi database global menggunakan Singleton pattern
func GetDB() *sql.DB {
	once.Do(func() {
		if DB != nil {
			return
		}
		var err error
		DB, err = ConnectDB()
		if err != nil {
			slog.Error("Failed to automatically connect to database in GetDB", "error", err)
		}
	})
	return DB
}
