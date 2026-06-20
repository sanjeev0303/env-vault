package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	MasterKey   string
	APIToken    string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, falling back to environment variables")
	}

	port := getEnv("PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "")
	masterKey := getEnv("MASTER_KEY", "")
	apiToken := getEnv("API_TOKEN", "")

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		MasterKey:   masterKey,
		APIToken:    apiToken,
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
