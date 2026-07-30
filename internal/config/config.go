package config

import (
	"os"
)

type Config struct {
	Port         string
	DatabasePath string
	JWTSecret    string
	AllowOrigins []string
}

func Load() Config {
	return Config{
		Port:         getEnv("PORT", "8080"),
		DatabasePath: getEnv("DATABASE_PATH", "coffee-site.db"),
		JWTSecret:    getEnv("JWT_SECRET", "dev-secret-change-me"),
		AllowOrigins: []string{
			getEnv("FRONTEND_ORIGIN", "http://localhost:5173"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
