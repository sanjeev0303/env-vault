package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	MasterKey   string
	APIToken    string // Legacy, kept for backward compat during migration

	// JWT
	JWTKeyDir       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration

	// Argon2
	Argon2Memory      uint32
	Argon2Iterations  uint32
	Argon2Parallelism uint8

	// CORS
	AllowedOrigins string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, falling back to environment variables")
	}

	port := getEnv("PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "")
	masterKey := getEnv("MASTER_KEY", "")
	apiToken := getEnv("API_TOKEN", "")
	jwtKeyDir := getEnv("JWT_KEY_DIR", "./keys")
	allowedOrigins := getEnv("ALLOWED_ORIGINS", "http://localhost:3000")

	return &Config{
		Port:              port,
		DatabaseURL:       dbURL,
		MasterKey:         masterKey,
		APIToken:          apiToken,
		JWTKeyDir:         jwtKeyDir,
		AccessTokenTTL:    15 * time.Minute,
		RefreshTokenTTL:   30 * 24 * time.Hour,
		Argon2Memory:      64 * 1024,
		Argon2Iterations:  3,
		Argon2Parallelism: 2,
		AllowedOrigins:    allowedOrigins,
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
