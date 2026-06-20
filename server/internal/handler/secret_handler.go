package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, err)
		return
	}

	sec, err := h.service.CreateSecret(r.Context(), projectID, envID, req.Key, req.Value)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, mapSecretToResp(sec))
}

func (h *SecretHandler) ListSecrets(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	envID := chi.URLParam(r, "envID")

	secrets, err := h.service.ListSecrets(r.Context(), projectID, envID)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	formatParam := r.URL.Query().Get("format")
	acceptHeader := r.Header.Get("Accept")

	if formatParam == "env" || strings.Contains(acceptHeader, "text/plain") {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		for _, s := range secrets {
			val := s.Value
			if strings.ContainsAny(val, " \t\n\r\"'\\#") {
				escaped := strings.ReplaceAll(val, "\\", "\\\\")
				escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
				escaped = strings.ReplaceAll(escaped, "\n", "\\n")
				escaped = strings.ReplaceAll(escaped, "\r", "\\r")
				val = fmt.Sprintf("\"%s\"", escaped)
			}
			_, _ = fmt.Fprintf(w, "%s=%s\n", s.Key, val)
		}
		return
	}

	resps := make([]dto.SecretResp, len(secrets))
	for i, s := range secrets {
		resps[i] = mapSecretToResp(s)
	}

	RespondWithJSON(w, http.StatusOK, resps)
}

func (h *SecretHandler) UpdateSecret(w http.ResponseWriter, r *http.Request) {
	secretID := chi.URLParam(r, "secretID")

	var req dto.UpdateSecretReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, err)
		return
	}

	sec, err := h.service.UpdateSecret(r.Context(), secretID, req.Value)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, mapSecretToResp(sec))
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

func mapSecretToResp(s *domain.Secret) dto.SecretResp {
	return dto.SecretResp{
		ID:            s.ID,
		ProjectID:     s.ProjectID,
		EnvironmentID: s.EnvironmentID,
		Key:           s.Key,
		Value:         s.Value,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}
