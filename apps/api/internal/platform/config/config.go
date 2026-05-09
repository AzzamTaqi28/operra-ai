package config

import "os"

type Config struct {
	AppEnv           string
	ListenAddr       string
	DatabaseURL      string
	JWTSecret        string
	APIURL           string
	AppURL           string
	StorageDriver    string
	S3Endpoint       string
	S3Bucket         string
	S3AccessKey      string
	S3SecretKey      string
	S3Region         string
	S3ForcePathStyle bool
	AIProvider       string
	AIBaseURL        string
	AIAPIKey         string
	AIModel          string
}

func Load() Config {
	return Config{
		AppEnv:           getEnv("APP_ENV", "development"),
		ListenAddr:       getEnv("API_LISTEN_ADDR", ":8080"),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://operra:operra@localhost:5432/operra?sslmode=disable"),
		JWTSecret:        getEnv("JWT_SECRET", "change-me-please"),
		APIURL:           getEnv("API_URL", "http://localhost:8080"),
		AppURL:           getEnv("APP_URL", "http://localhost:3000"),
		StorageDriver:    getEnv("STORAGE_DRIVER", "local"),
		S3Endpoint:       getEnv("S3_ENDPOINT", ""),
		S3Bucket:         getEnv("S3_BUCKET", "operra"),
		S3AccessKey:      getEnv("S3_ACCESS_KEY", ""),
		S3SecretKey:      getEnv("S3_SECRET_KEY", ""),
		S3Region:         getEnv("S3_REGION", "us-east-1"),
		S3ForcePathStyle: getEnv("S3_FORCE_PATH_STYLE", "true") == "true",
		AIProvider:       getEnv("AI_PROVIDER", "openai"),
		AIBaseURL:        getEnv("AI_BASE_URL", "https://api.openai.com/v1"),
		AIAPIKey:         getEnv("AI_API_KEY", ""),
		AIModel:          getEnv("AI_MODEL", ""),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
