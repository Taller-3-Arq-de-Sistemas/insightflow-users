package config

import (
	"os"
)

type Config struct {
	Port      string
	DBUrl     string
	JWTSecret string
}

func GetString(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func Load() Config {
	return Config{
		Port:      GetString("PORT", "8080"),
		DBUrl:     GetString("DB_URL", "postgres://postgres:postgres@localhost:5432/insightflow_users?sslmode=disable"),
		JWTSecret: GetString("JWT_SECRET", "secret"),
	}
}
