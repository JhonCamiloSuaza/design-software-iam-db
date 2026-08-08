package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
	JWTSecret   string
	FrontendURL string
}

func Load() Config {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		env("DB_USER", "iam_user"),
		env("DB_PASSWORD", "iam_pass"),
		env("DB_HOST", "localhost"),
		env("DB_PORT", "5436"),
		env("DB_NAME", "iam_db"),
	)
	return Config{
		DatabaseURL: databaseURL,
		Port:        env("PORT", "8080"),
		JWTSecret:   env("JWT_SECRET", "iam_demo_secret_change_in_production"),
		FrontendURL: env("FRONTEND_URL", "http://localhost:5173"),
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
