package router

import (
	"github.com/go-chi/chi/v5"

	"env-vault/server/internal/handler"
)

func RegisterAuthRoutes(r chi.Router, h *handler.AuthHandler) {
	r.Post("/auth/register", h.Register)
	r.Post("/auth/login", h.Login)
	r.Post("/auth/refresh", h.Refresh)
	r.Get("/auth/csrf", h.GetCSRFToken)
	r.Post("/auth/forgot-password", h.ForgotPassword)
	r.Post("/auth/reset-password", h.ResetPassword)
}

func RegisterProtectedAuthRoutes(r chi.Router, h *handler.AuthHandler) {
	r.Post("/auth/logout", h.Logout)
	r.Get("/auth/me", h.Me)
	r.Post("/auth/reauth", h.Reauth)
}
