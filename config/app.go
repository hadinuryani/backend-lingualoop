package config

import (
	"strconv"
)

type App struct {
	Name  string
	Env   string
	Port  string
	Debug bool
}

func LoadAppConfig() App {
	debugMode, _ := strconv.ParseBool(getEnv("APP_DEBUG", "false"))
	return App{
		Name:  getEnv("APP_NAME", "backend-lingualoop"),
		Env:   getEnv("APP_ENV", "development"),
		Port:  getEnv("APP_PORT", "8080"),
		Debug: debugMode,
	}
}

