package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration
type Config struct {
	Port              string
	DatabaseURL       string
	PostgresSchema    string
	GoogleClientID    string
	GoogleClientSecret string
	GoogleRedirectURL string
	WorkspaceDomain   string
	JWTSecret         string
	FrontendURL       string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Port:              getEnv("PORT", "8080"),
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		PostgresSchema:    getEnv("POSTGRES_SCHEMA", "public"),
		GoogleClientID:    getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL: getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/callback"),
		WorkspaceDomain:   getEnv("WORKSPACE_DOMAIN", ""),
		JWTSecret:         getEnv("JWT_SECRET", ""),
		FrontendURL:       getEnv("FRONTEND_URL", "http://localhost:3000"),
	}

	// Validate required fields
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
