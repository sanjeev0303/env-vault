package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"env-vault/server/internal/domain"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
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
	message := "an unexpected error occurred"
	code := "INTERNAL_ERROR"

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			status = http.StatusNotFound
			message = "resource not found"
			code = "NOT_FOUND"
		case errors.Is(err, domain.ErrAlreadyExists):
			status = http.StatusConflict
			message = "resource already exists"
			code = "ALREADY_EXISTS"
		case errors.Is(err, domain.ErrInvalidInput):
			status = http.StatusBadRequest
			message = "invalid input data"
			code = "INVALID_INPUT"
		case errors.Is(err, domain.ErrUnauthorized):
			status = http.StatusUnauthorized
			message = "unauthorized"
			code = "UNAUTHORIZED"
		case errors.Is(err, domain.ErrForbidden):
			status = http.StatusForbidden
			message = "access denied"
			code = "FORBIDDEN"
		default:
			// Log internal error but don't expose details
			log.Printf("internal error: %v", err)
		}
	}

	RespondWithJSON(w, status, ErrorResponse{Error: message, Code: code})
}
