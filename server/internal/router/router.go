package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"env-vault/server/internal/domain"
	"env-vault/server/internal/handler"
)

func NewRouter(
	projectHandler *handler.ProjectHandler,
	secretHandler *handler.SecretHandler,
	authHandler *handler.AuthHandler,
	mw *handler.Middleware,
	allowedOrigins string,
) *chi.Mux {
	r := chi.NewRouter()

	apiLimiter := handler.NewRateLimiter(10, 50)  // 10 req/sec, burst 50
	authLimiter := handler.NewRateLimiter(2, 10)  // 2 req/sec, burst 10

	// Core Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Security Headers
	r.Use(securityHeaders)

	// CORS Middleware
	r.Use(corsMiddleware(allowedOrigins))

	// Public Health Check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Protected API Routes (JWT required)
	r.Route("/api", func(api chi.Router) {
		api.Use(apiLimiter.Limit)
		
		// Public Auth Routes (no JWT required)
		api.Group(func(publicAuth chi.Router) {
			publicAuth.Use(authLimiter.Limit)
			RegisterAuthRoutes(publicAuth, authHandler)
		})

		// Protected endpoints
		api.Group(func(protected chi.Router) {
			protected.Use(mw.RequireJWT)
			protected.Use(mw.RequireCSRF)

			protected.Get("/metrics", handler.MetricsHandler)

			// Protected auth routes
			RegisterProtectedAuthRoutes(protected, authHandler)

			// Project routes with RBAC
			protected.Post("/projects", projectHandler.CreateProject)
			protected.Get("/projects", projectHandler.ListProjects)

			protected.Route("/projects/{projectID}", func(pr chi.Router) {
				pr.Use(mw.RequirePermission(domain.PermViewSecret))
				pr.Get("/", projectHandler.GetProject)
				pr.Get("/environments", projectHandler.ListEnvironments)

				pr.Group(func(manage chi.Router) {
					manage.Use(mw.RequirePermission(domain.PermManageProject))
					manage.Delete("/", projectHandler.DeleteProject)
				})

				// Environment management
				pr.Group(func(env chi.Router) {
					env.Use(mw.RequirePermission(domain.PermManageProject))
					env.Post("/environments", projectHandler.CreateEnvironment)
					env.Delete("/environments/{envID}", projectHandler.DeleteEnvironment)
				})

				// Secret routes
				pr.Route("/environments/{envID}/secrets", func(sr chi.Router) {
					sr.With(mw.RequirePermission(domain.PermViewSecret)).Get("/", secretHandler.ListSecrets)
					sr.With(mw.RequirePermission(domain.PermViewSecret), mw.RequireReauth).Get("/export", secretHandler.ExportSecrets)

					sr.With(mw.RequirePermission(domain.PermCreateSecret)).Post("/", secretHandler.CreateSecret)

					sr.Route("/{secretID}", func(secr chi.Router) {
						secr.With(mw.RequirePermission(domain.PermViewSecret)).Get("/reveal", secretHandler.RevealSecret)
						secr.With(mw.RequirePermission(domain.PermViewSecret)).Get("/history", secretHandler.GetSecretHistory)
						secr.With(mw.RequirePermission(domain.PermUpdateSecret)).Put("/", secretHandler.UpdateSecret)
						secr.With(mw.RequirePermission(domain.PermDeleteSecret)).Delete("/", secretHandler.DeleteSecret)
						secr.With(mw.RequirePermission(domain.PermUpdateSecret)).Post("/rollback", secretHandler.RollbackSecret)
					})
				})

				// Members
				pr.Route("/members", func(mr chi.Router) {
					mr.Use(mw.RequirePermission(domain.PermManageMembers))
					mr.Get("/", projectHandler.ListMembers)
					mr.Post("/", projectHandler.AddMember)
					mr.Put("/{memberUserID}", projectHandler.UpdateMemberRole)
					mr.Delete("/{memberUserID}", projectHandler.RemoveMember)
				})
			})
		})
	})

	return r
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none'")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(allowedOrigins string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-CSRF-Token")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
