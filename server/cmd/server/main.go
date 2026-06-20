package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"env-vault/server/internal/config"
	"env-vault/server/internal/database"
	"env-vault/server/internal/encryption"
	"env-vault/server/internal/handler"
	"env-vault/server/internal/repository/postgres"
	"env-vault/server/internal/router"
	"env-vault/server/internal/service"
)

func main() {
	log.Println("Starting env-vault server...")

	// 1. Load config
	cfg := config.LoadConfig()

	// 2. Initialize Database
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Close()
	log.Println("Connected to PostgreSQL database")

	// 3. Initialize Encryption Service
	encryptor, err := encryption.NewAESEncryptionService(cfg.MasterKey)
	if err != nil {
		log.Fatalf("Encryption service initialization failed: %v", err)
	}
	log.Println("Encryption service initialized successfully")

	// 4. Initialize Repositories
	projectRepo := postgres.NewProjectRepository(db)
	secretRepo := postgres.NewSecretRepository(db)

	// 5. Initialize Services
	projectSvc := service.NewProjectService(projectRepo)
	secretSvc := service.NewSecretService(secretRepo, projectRepo, encryptor)

	// 6. Initialize Handlers
	projectHandler := handler.NewProjectHandler(projectSvc)
	secretHandler := handler.NewSecretHandler(secretSvc)

	// 7. Setup Router
	r := router.NewRouter(projectHandler, secretHandler, cfg.APIToken)

	// 8. Start HTTP Server with Graceful Shutdown
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server listening on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server listen failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
