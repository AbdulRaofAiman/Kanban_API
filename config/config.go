package config

import (
	"os"
)

type AppConfig struct {
	Port      string
	JWTSecret string
	JWTExpiry string
	AppEnv    string
}

func GetConfig() AppConfig {
	return AppConfig{
		Port:      os.Getenv("PORT"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		JWTExpiry: os.Getenv("JWT_EXPIRY"),
		AppEnv:    os.Getenv("APP_ENV"),
	}
}
