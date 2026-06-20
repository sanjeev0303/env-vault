package handler

import (
	"net/http"
	"strings"

	"env-vault/server/internal/domain"
)

type Middleware struct {
	apiToken string
}

func NewMiddleware(apiToken string) *Middleware {
	return &Middleware{apiToken: apiToken}
}

// RequireAuth protects routes with a bearer token.
func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.apiToken == "" {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			RespondWithError(w, domain.ErrUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			RespondWithError(w, domain.ErrUnauthorized)
			return
		}

		token := parts[1]
		if token != m.apiToken {
			RespondWithError(w, domain.ErrUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
