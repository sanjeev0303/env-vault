package handler

import (
	"net/http"

	"env-vault/server/internal/dto"
	"env-vault/server/internal/service"

	"github.com/google/uuid"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterReq
	if err := ParseAndValidate(r, &req); err != nil {
		RespondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.service.Register(r.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, dto.UserResp{
		ID:           user.ID,
		Email:        user.Email,
		Name:         user.Name,
		IsSuperAdmin: user.IsSuperAdmin,
		CreatedAt:    user.CreatedAt,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginReq
	if err := ParseAndValidate(r, &req); err != nil {
		RespondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	ip := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = forwarded
	}

	user, _, accessToken, refreshToken, err := h.service.Login(r.Context(), req.Email, req.Password, ip, r.UserAgent())
	if err != nil {
		RespondWithError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/auth/refresh",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   7 * 24 * 3600,
	})

	RespondWithJSON(w, http.StatusOK, dto.AuthResp{
		AccessToken:  accessToken,
		ExpiresIn:    900, // 15 minutes
		User: dto.UserResp{
			ID:           user.ID,
			Email:        user.Email,
			Name:         user.Name,
			IsSuperAdmin: user.IsSuperAdmin,
			CreatedAt:    user.CreatedAt,
		},
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID := GetSessionIDFromContext(r.Context())
	
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/auth/refresh",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})

	if sessionID == "" {
		RespondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	if err := h.service.Logout(r.Context(), sessionID); err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		RespondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: "missing refresh token cookie"})
		return
	}

	accessToken, newRefreshToken, err := h.service.RefreshToken(r.Context(), cookie.Value)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/api/auth/refresh",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   7 * 24 * 3600,
	})

	RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"access_token":  accessToken,
		"expires_in":    900,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		RespondWithError(w, ErrUnauthorizedAccess)
		return
	}

	user, err := h.service.GetCurrentUser(r.Context(), userID)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, dto.UserResp{
		ID:           user.ID,
		Email:        user.Email,
		Name:         user.Name,
		IsSuperAdmin: user.IsSuperAdmin,
		CreatedAt:    user.CreatedAt,
	})
}

// GetCSRFToken issues a new CSRF token to the client.
func (h *AuthHandler) GetCSRFToken(w http.ResponseWriter, r *http.Request) {
	token := uuid.NewString()

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   3600, // 1 hour
	})

	RespondWithJSON(w, http.StatusOK, map[string]string{"csrf_token": token})
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ForgotPasswordReq
	if err := ParseAndValidate(r, &req); err != nil {
		RespondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// In a real system, we'd throttle this by IP or email
	err := h.service.ForgotPassword(r.Context(), req.Email)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "If an account with that email exists, a reset link has been sent."})
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ResetPasswordReq
	if err := ParseAndValidate(r, &req); err != nil {
		RespondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.service.ResetPassword(r.Context(), req.Token, req.NewPassword)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Password successfully reset."})
}

func (h *AuthHandler) Reauth(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		RespondWithError(w, ErrUnauthorizedAccess)
		return
	}

	var req dto.ReauthReq
	if err := ParseAndValidate(r, &req); err != nil {
		RespondWithJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	reauthToken, err := h.service.Reauth(r.Context(), userID, req.Password)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"reauth_token": reauthToken})
}
