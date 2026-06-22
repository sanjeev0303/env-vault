package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"env-vault/server/internal/auth"
	"env-vault/server/internal/domain"
	"env-vault/server/internal/repository"

	"github.com/go-chi/chi/v5"
)



var ErrUnauthorizedAccess = domain.ErrUnauthorized

type Middleware struct {
	apiToken   string // Legacy support
	jwtService *auth.JWTService
	memberRepo repository.ProjectMemberRepository
}

func NewMiddleware(apiToken string, jwtService *auth.JWTService, memberRepo repository.ProjectMemberRepository) *Middleware {
	return &Middleware{
		apiToken:   apiToken,
		jwtService: jwtService,
		memberRepo: memberRepo,
	}
}

// RequireJWT validates JWT access tokens and injects user context.
func (m *Middleware) RequireJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			RespondWithError(w, domain.ErrUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			RespondWithError(w, domain.ErrUnauthorized)
			return
		}

		claims, err := m.jwtService.ValidateAccessToken(parts[1])
		if err != nil {
			RespondWithError(w, domain.ErrUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, domain.CtxUserID, claims.UserID)
		ctx = context.WithValue(ctx, domain.CtxUserEmail, claims.Email)
		ctx = context.WithValue(ctx, domain.CtxSessionID, claims.SessionID)
		ctx = context.WithValue(ctx, domain.CtxIsSuperAdmin, claims.IsSuperAdmin)

		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = forwarded
		}
		ctx = context.WithValue(ctx, domain.CtxIPAddress, ip)
		ctx = context.WithValue(ctx, domain.CtxUserAgent, r.UserAgent())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequirePermission checks that the authenticated user has the given permission on the project from the URL.
func (m *Middleware) RequirePermission(perm domain.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := GetUserIDFromContext(r.Context())
			if userID == "" {
				RespondWithError(w, domain.ErrUnauthorized)
				return
			}

			// Superadmins bypass RBAC
			if GetIsSuperAdminFromContext(r.Context()) {
				next.ServeHTTP(w, r)
				return
			}

			projectID := chi.URLParam(r, "projectID")
			if projectID == "" {
				RespondWithError(w, domain.ErrInvalidInput)
				return
			}

			member, err := m.memberRepo.GetMember(r.Context(), projectID, userID)
			if err != nil {
				if errors.Is(err, domain.ErrNotFound) {
					RespondWithJSON(w, http.StatusForbidden, ErrorResponse{Error: "access denied"})
					return
				}
				RespondWithError(w, err)
				return
			}

			if !domain.HasPermission(member.Role, perm) {
				RespondWithJSON(w, http.StatusForbidden, ErrorResponse{Error: "insufficient permissions"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireSuperAdmin restricts access to superadmin users only.
func (m *Middleware) RequireSuperAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !GetIsSuperAdminFromContext(r.Context()) {
			RespondWithJSON(w, http.StatusForbidden, ErrorResponse{Error: "superadmin access required"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Context helpers
func GetUserIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(domain.CtxUserID).(string)
	return v
}

func GetSessionIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(domain.CtxSessionID).(string)
	return v
}

func GetIsSuperAdminFromContext(ctx context.Context) bool {
	v, _ := ctx.Value(domain.CtxIsSuperAdmin).(bool)
	return v
}

func GetIPAddressFromContext(ctx context.Context) string {
	v, _ := ctx.Value(domain.CtxIPAddress).(string)
	return v
}

func GetUserAgentFromContext(ctx context.Context) string {
	v, _ := ctx.Value(domain.CtxUserAgent).(string)
	return v
}

// RequireCSRF enforces Double Submit Cookie CSRF validation for mutating requests.
func (m *Middleware) RequireCSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only enforce on mutating methods
		if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("csrf_token")
		if err != nil || cookie.Value == "" {
			RespondWithJSON(w, http.StatusForbidden, ErrorResponse{Error: "missing csrf cookie"})
			return
		}

		headerToken := r.Header.Get("X-CSRF-Token")
		if headerToken == "" || headerToken != cookie.Value {
			RespondWithJSON(w, http.StatusForbidden, ErrorResponse{Error: "invalid csrf token"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequireReauth validates a short-lived reauth token for sensitive actions.
func (m *Middleware) RequireReauth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("X-Reauth-Token")
		if tokenStr == "" {
			RespondWithError(w, ErrUnauthorizedAccess)
			return
		}

		claims, err := m.jwtService.ValidateAccessToken(tokenStr)
		if err != nil || !claims.IsReauth {
			RespondWithError(w, ErrUnauthorizedAccess)
			return
		}

		next.ServeHTTP(w, r)
	})
}
