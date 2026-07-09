package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"backend-lingualoop/config"

	_ "github.com/go-sql-driver/mysql"
)

// DB adalah instance koneksi database global
var DB *sql.DB

func ConnectDB() *sql.DB {
	cfg := config.GetConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Gagal membuka koneksi database: %v", err)
	}

	// Konfigurasi Connection Pool
	DB.SetMaxOpenConns(cfg.Database.MaxConnections)
	DB.SetMaxIdleConns(cfg.Database.MaxIdleConnections)
	DB.SetConnMaxLifetime(time.Hour)
	DB.SetConnMaxIdleTime(15 * time.Minute)

	// Ping ke database untuk memastikan koneksi berhasil
	if err := DB.Ping(); err != nil {
		log.Fatalf("Gagal terhubung ke database (Ping): %v", err)
	}

	log.Println("Koneksi ke database MySQL berhasil dibentuk")
	return DB
}

// GetDB mengembalikan instance koneksi database global
func GetDB() *sql.DB {
	if DB == nil {
		return ConnectDB()
	}
	return DB
}
