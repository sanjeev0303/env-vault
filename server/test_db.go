package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file, continuing with env vars")
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer db.Close()

	rows, err := db.QueryContext(context.Background(), "SELECT email, password_hash FROM users")
	if err != nil {
		log.Fatal("Query failed:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var email, hash string
		if err := rows.Scan(&email, &hash); err != nil {
			log.Fatal("Scan failed:", err)
		}
		fmt.Printf("Email: %s\nHash: %s\n\n", email, hash)
	}
}
