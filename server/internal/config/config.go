package config

import (
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	MasterKey   string
	APIToken    string
}

func LoadConfig() *Config {
	port := getEnv("PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "postgresql://neondb_owner:npg_g23KeNMfLqwh@ep-quiet-tree-atbnsxzh-pooler.c-9.us-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require")
	masterKey := getEnv("MASTER_KEY", "dev-master-key-change-me-in-production-123456")
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
