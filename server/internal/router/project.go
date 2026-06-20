package router

import (
	"github.com/go-chi/chi/v5"

	"env-vault/server/internal/handler"
)

func RegisterProjectRoutes(r chi.Router, h *handler.ProjectHandler) {
	r.Post("/projects", h.CreateProject)
	r.Get("/projects", h.ListProjects)
	r.Get("/projects/{projectID}", h.GetProject)
	r.Delete("/projects/{projectID}", h.DeleteProject)

	r.Post("/projects/{projectID}/environments", h.CreateEnvironment)
	r.Get("/projects/{projectID}/environments", h.ListEnvironments)
	r.Delete("/projects/{projectID}/environments/{envID}", h.DeleteEnvironment)
}
