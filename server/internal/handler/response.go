package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"env-vault/server/internal/domain"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

func RespondWithError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	message := "internal server error"

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			status = http.StatusNotFound
			message = err.Error()
		case errors.Is(err, domain.ErrAlreadyExists):
			status = http.StatusConflict
			message = err.Error()
		case errors.Is(err, domain.ErrInvalidInput):
			status = http.StatusBadRequest
			message = err.Error()
		case errors.Is(err, domain.ErrUnauthorized):
			status = http.StatusUnauthorized
			message = err.Error()
		default:
			message = err.Error()
		}
	}

	RespondWithJSON(w, status, ErrorResponse{Error: message})
}
