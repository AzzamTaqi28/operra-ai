package config

import "os"

type Config struct {
	AppEnv      string
	ListenAddr  string
	DatabaseURL string
	JWTSecret   string
	APIURL      string
	AppURL      string
}

func Load() Config {
	return Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		ListenAddr:  getEnv("API_LISTEN_ADDR", ":8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://operra:operra@localhost:5432/operra?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me-please"),
		APIURL:      getEnv("API_URL", "http://localhost:8080"),
		AppURL:      getEnv("APP_URL", "http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
