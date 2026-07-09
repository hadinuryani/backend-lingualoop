package config

import "strconv"

type Database struct {
	Host               string
	Port               string
	User               string
	Password           string
	Name               string
	MaxConnections     int
	MaxIdleConnections int
}

func LoadDatabaseConfig() Database {
	maxConn, _ := strconv.Atoi(getEnv("DB_MAX_CONNECTIONS", "100"))
	maxIdleConn, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNECTIONS", "10"))
	return Database{
		Host:               getEnv("DB_HOST", "localhost"),
		Port:               getEnv("DB_PORT", "3306"),
		User:               getEnv("DB_USER", "root"),
		Password:           getEnv("DB_PASSWORD", ""),
		Name:               getEnv("DB_NAME", "lingualoop"),
		MaxConnections:     maxConn,
		MaxIdleConnections: maxIdleConn,
	}
}

