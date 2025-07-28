 package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application.
type Config struct {
	AppPort string
	DSN     string
	JWT_SECRET string
	ADMIN_EMAIL string
	ADMIN_PASSWORD string
	EMAIL_API_KEY string
	EMAIL_SENDER string
}

// New creates and returns a new Config struct, populated from environment variables.
func New() *Config {
	// Load .env file. It's safe to ignore the error if it doesn't exist,
	// as env vars can be set directly in a deployment environment.
	if err := godotenv.Load(".env"); err != nil {
		log.Println("INFO: No .env file found, using system environment variables or defaults.")
	}

	return &Config{
		AppPort: getEnv("APP_PORT", "3001"),
		DSN:     getEnv("DATABASE_DSN", ""),
		JWT_SECRET : getEnv ("JWT_SECRET","secret-key"),
		ADMIN_EMAIL: getEnv("ADMIN_EMAIL", ""),
		ADMIN_PASSWORD: getEnv("ADMIN_PASSWORD", ""),
		EMAIL_API_KEY: getEnv("EMAIL_API_KEY", ""),
		EMAIL_SENDER: getEnv("EMAIL_SENDER", ""),
	}
}

// getEnv retrieves an environment variable or returns a fallback value.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}


