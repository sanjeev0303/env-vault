package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"env-vault/server/internal/domain"
	"env-vault/server/internal/dto"
	"env-vault/server/internal/service"
)

type ProjectHandler struct {
	service service.ProjectService
}

func NewProjectHandler(service service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProjectReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, err)
		return
	}

	p, err := h.service.CreateProject(r.Context(), req.Name)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, mapProjectToResp(p))
}

func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.service.ListProjects(r.Context())
	if err != nil {
		RespondWithError(w, err)
		return
	}

	resps := make([]dto.ProjectResp, len(projects))
	for i, p := range projects {
		resps[i] = mapProjectToResp(p)
	}

	RespondWithJSON(w, http.StatusOK, resps)
}

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "projectID")
	p, err := h.service.GetProject(r.Context(), id)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, mapProjectToResp(p))
}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "projectID")
	err := h.service.DeleteProject(r.Context(), id)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}

func (h *ProjectHandler) CreateEnvironment(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	var req dto.CreateEnvReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, err)
		return
	}

	env, err := h.service.CreateEnvironment(r.Context(), projectID, req.Name)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, mapEnvToResp(env))
}

func (h *ProjectHandler) ListEnvironments(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	envs, err := h.service.ListEnvironments(r.Context(), projectID)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	resps := make([]dto.EnvResp, len(envs))
	for i, env := range envs {
		resps[i] = mapEnvToResp(env)
	}

	RespondWithJSON(w, http.StatusOK, resps)
}

func (h *ProjectHandler) DeleteEnvironment(w http.ResponseWriter, r *http.Request) {
	envID := chi.URLParam(r, "envID")
	err := h.service.DeleteEnvironment(r.Context(), envID)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}

func mapProjectToResp(p *domain.Project) dto.ProjectResp {
	return dto.ProjectResp{
		ID:        p.ID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func mapEnvToResp(env *domain.Environment) dto.EnvResp {
	return dto.EnvResp{
		ID:        env.ID,
		ProjectID: env.ProjectID,
		Name:      env.Name,
		CreatedAt: env.CreatedAt,
		UpdatedAt: env.UpdatedAt,
	}
}
