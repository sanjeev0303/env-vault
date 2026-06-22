package handler

import (
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
	if err := ParseAndValidate(r, &req); err != nil {
		RespondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userID := GetUserIDFromContext(r.Context())
	p, err := h.service.CreateProject(r.Context(), req.Name, userID)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, mapProjectToResp(p))
}

func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	isSuperAdmin := GetIsSuperAdminFromContext(r.Context())

	projects, err := h.service.ListProjects(r.Context(), userID, isSuperAdmin)
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
	if err := ParseAndValidate(r, &req); err != nil {
		RespondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
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

// Member management handlers
func (h *ProjectHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	members, err := h.service.ListMembers(r.Context(), projectID)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	resps := make([]dto.MemberResp, len(members))
	for i, m := range members {
		resps[i] = dto.MemberResp{
			ID:        m.ID,
			ProjectID: m.ProjectID,
			UserID:    m.UserID,
			Role:      string(m.Role),
			CreatedAt: m.CreatedAt,
		}
	}

	RespondWithJSON(w, http.StatusOK, resps)
}

func (h *ProjectHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	var req dto.AddMemberReq
	if err := ParseAndValidate(r, &req); err != nil {
		RespondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	member, err := h.service.AddMember(r.Context(), projectID, req.Email, domain.Role(req.Role))
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, dto.MemberResp{
		ID:        member.ID,
		ProjectID: member.ProjectID,
		UserID:    member.UserID,
		Role:      string(member.Role),
		CreatedAt: member.CreatedAt,
	})
}

func (h *ProjectHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	memberUserID := chi.URLParam(r, "memberUserID")
	var req dto.UpdateMemberReq
	if err := ParseAndValidate(r, &req); err != nil {
		RespondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.service.UpdateMemberRole(r.Context(), projectID, memberUserID, domain.Role(req.Role))
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *ProjectHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	memberUserID := chi.URLParam(r, "memberUserID")

	err := h.service.RemoveMember(r.Context(), projectID, memberUserID)
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
