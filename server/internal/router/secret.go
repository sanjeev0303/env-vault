package router

import (
	"github.com/go-chi/chi/v5"

	"env-vault/server/internal/handler"
)

func RegisterSecretRoutes(r chi.Router, h *handler.SecretHandler) {
	r.Post("/projects/{projectID}/environments/{envID}/secrets", h.CreateSecret)
	r.Get("/projects/{projectID}/environments/{envID}/secrets", h.ListSecrets)
	r.Put("/projects/{projectID}/environments/{envID}/secrets/{secretID}", h.UpdateSecret)
	r.Delete("/projects/{projectID}/environments/{envID}/secrets/{secretID}", h.DeleteSecret)
}
