package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"env-vault/server/internal/domain"
	"env-vault/server/internal/dto"
	"env-vault/server/internal/service"
)

type SecretHandler struct {
	service service.SecretService
}

func NewSecretHandler(service service.SecretService) *SecretHandler {
	return &SecretHandler{service: service}
}

func (h *SecretHandler) CreateSecret(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	envID := chi.URLParam(r, "envID")

	var req dto.CreateSecretReq
	if err := ParseAndValidate(r, &req); err != nil {
		RespondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	sec, err := h.service.CreateSecret(r.Context(), projectID, envID, req.Key, req.Value)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, mapSecretToMetadata(sec))
}

// ListSecrets returns metadata only — no secret values.
func (h *SecretHandler) ListSecrets(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	envID := chi.URLParam(r, "envID")

	cursor := r.URL.Query().Get("cursor")
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	secrets, nextCursor, err := h.service.ListSecrets(r.Context(), projectID, envID, cursor, limit)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	resps := make([]dto.SecretMetadataResp, len(secrets))
	for i, s := range secrets {
		resps[i] = mapSecretToMetadata(s)
	}

	// Returning a paginated structure
	RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data":        resps,
		"next_cursor": nextCursor,
	})
}

// RevealSecret decrypts and returns a single secret value.
func (h *SecretHandler) RevealSecret(w http.ResponseWriter, r *http.Request) {
	secretID := chi.URLParam(r, "secretID")

	sec, err := h.service.RevealSecret(r.Context(), secretID)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, dto.SecretRevealResp{
		ID:            sec.ID,
		ProjectID:     sec.ProjectID,
		EnvironmentID: sec.EnvironmentID,
		Key:           sec.Key,
		Value:         sec.Value,
		Version:       sec.Version,
		CreatedAt:     sec.CreatedAt,
		UpdatedAt:     sec.UpdatedAt,
	})
}

// ExportSecrets exports all secrets for an environment (requires reauth)
func (h *SecretHandler) ExportSecrets(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	envID := chi.URLParam(r, "envID")

	secretsMap, err := h.service.ExportSecrets(r.Context(), projectID, envID)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, secretsMap)
}

func (h *SecretHandler) UpdateSecret(w http.ResponseWriter, r *http.Request) {
	secretID := chi.URLParam(r, "secretID")

	var req dto.UpdateSecretReq
	if err := ParseAndValidate(r, &req); err != nil {
		RespondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	sec, err := h.service.UpdateSecret(r.Context(), secretID, req.Value)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, mapSecretToMetadata(sec))
}

func (h *SecretHandler) DeleteSecret(w http.ResponseWriter, r *http.Request) {
	secretID := chi.URLParam(r, "secretID")

	err := h.service.DeleteSecret(r.Context(), secretID)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}

func (h *SecretHandler) GetSecretHistory(w http.ResponseWriter, r *http.Request) {
	secretID := chi.URLParam(r, "secretID")

	versions, err := h.service.GetSecretHistory(r.Context(), secretID)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	resps := make([]dto.SecretVersionResp, len(versions))
	for i, v := range versions {
		resps[i] = dto.SecretVersionResp{
			ID:        v.ID,
			SecretID:  v.SecretID,
			Version:   v.Version,
			CreatedAt: v.CreatedAt,
		}
	}

	RespondWithJSON(w, http.StatusOK, resps)
}

func (h *SecretHandler) RollbackSecret(w http.ResponseWriter, r *http.Request) {
	secretID := chi.URLParam(r, "secretID")

	var req dto.RollbackReq
	if err := ParseAndValidate(r, &req); err != nil {
		RespondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	sec, err := h.service.RollbackSecret(r.Context(), secretID, req.Version)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, mapSecretToMetadata(sec))
}

func mapSecretToMetadata(s *domain.Secret) dto.SecretMetadataResp {
	return dto.SecretMetadataResp{
		ID:            s.ID,
		ProjectID:     s.ProjectID,
		EnvironmentID: s.EnvironmentID,
		Key:           s.Key,
		Version:       s.Version,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}
