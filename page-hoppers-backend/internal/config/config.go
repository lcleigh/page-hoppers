package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for our application
type Config struct {
	// Server configuration
	Port string
	Host string

	// Database configuration
	DatabaseURL string

	// JWT configuration
	JWTSecret string

	// Environment
	Environment string
	InDocker    bool

	// CORS configuration
	CORSAllowedOrigins []string
	CORSAllowedHeaders []string
	CORSAllowedMethods []string
}

// Load reads configuration from environment variables
func Load() *Config {
	// Load .env file if not in Docker
	if os.Getenv("IN_DOCKER") != "true" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found (this is fine if running in Docker)")
		}
	}

	config := &Config{
		// Server defaults
		Port: getEnv("PORT", "8080"),
		Host: getEnv("HOST", "localhost"),

		// Database
		DatabaseURL: getEnv("DATABASE_URL", ""),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", ""),

		// Environment
		Environment: getEnv("ENVIRONMENT", "development"),
		InDocker:    getEnvBool("IN_DOCKER", false),

		// CORS defaults
		CORSAllowedOrigins: getEnvSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
		CORSAllowedHeaders: getEnvSlice("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization"}),
		CORSAllowedMethods: getEnvSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
	}

	// Validate required configuration
	if config.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	if config.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	return config
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool gets a boolean environment variable with a fallback default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// getEnvSlice gets a slice environment variable with a fallback default value
func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Simple comma-separated values parsing
		// You could enhance this to handle more complex formats
		values := []string{}
		for _, v := range splitString(value, ",") {
			if trimmed := trimString(v); trimmed != "" {
				values = append(values, trimmed)
			}
		}
		if len(values) > 0 {
			return values
		}
	}
	return defaultValue
}

// Helper functions for string manipulation
func splitString(s, sep string) []string {
	result := []string{}
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trimString(s string) string {
	start := 0
	end := len(s)
	
	// Trim leading whitespace
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	
	// Trim trailing whitespace
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	
	return s[start:end]
}
