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

	"env-vault/server/internal/auth"
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

	// 3. Initialize JWT Service
	jwtService, err := auth.NewJWTService(cfg.JWTKeyDir, cfg.AccessTokenTTL)
	if err != nil {
		log.Fatalf("JWT service initialization failed: %v", err)
	}
	log.Println("JWT service initialized (RS256)")

	// 4. Initialize Encryption Service
	encryptor, err := encryption.NewAESEncryptionService(cfg.MasterKey)
	if err != nil {
		log.Fatalf("Encryption service initialization failed: %v", err)
	}
	log.Println("Encryption service initialized")

	// 5. Initialize Repositories
	userRepo := postgres.NewUserRepository(db)
	sessionRepo := postgres.NewSessionRepository(db)
	tokenRepo := postgres.NewRefreshTokenRepository(db)
	memberRepo := postgres.NewProjectMemberRepository(db)
	projectRepo := postgres.NewProjectRepository(db)
	secretRepo := postgres.NewSecretRepository(db)
	resetRepo := postgres.NewPasswordResetTokenRepository(db)
	auditRepo := postgres.NewAuditRepository(db)

	// 6. Initialize Services
	argon2Cfg := auth.Argon2Config{
		Memory:      cfg.Argon2Memory,
		Iterations:  cfg.Argon2Iterations,
		Parallelism: cfg.Argon2Parallelism,
		SaltLength:  16,
		KeyLength:   32,
	}
	authSvc := service.NewAuthService(userRepo, sessionRepo, tokenRepo, resetRepo, jwtService, argon2Cfg, cfg.RefreshTokenTTL)
	projectSvc := service.NewProjectService(projectRepo, memberRepo, userRepo)
	auditSvc := service.NewAuditService(auditRepo)
	secretSvc := service.NewSecretService(secretRepo, projectRepo, encryptor, auditSvc)

	// 7. Initialize Handlers
	authHandler := handler.NewAuthHandler(authSvc)
	projectHandler := handler.NewProjectHandler(projectSvc)
	secretHandler := handler.NewSecretHandler(secretSvc)

	// 8. Initialize Middleware
	mw := handler.NewMiddleware(cfg.APIToken, jwtService, memberRepo)

	// 9. Setup Router
	r := router.NewRouter(projectHandler, secretHandler, authHandler, mw, cfg.AllowedOrigins)

	// 10. Start HTTP Server with Graceful Shutdown
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
