package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	App      App
	Database Database
	JWT      JWT
	Gmail    Gmail
}

var cfg *Config

// LoadConfig membaca file .env dan mengembalikan struktur Config
func LoadConfig() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	cfg = &Config{
		App:      LoadAppConfig(),
		Database: LoadDatabaseConfig(),
		JWT:      LoadJWTConfig(),
		Gmail:    LoadGmailConfig(),
	}

	return cfg
}

// GetConfig mengembalikan konfigurasi yang sudah dimuat
func GetConfig() *Config {
	if cfg == nil {
		return LoadConfig()
	}
	return cfg
}
