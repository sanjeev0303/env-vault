package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"env-vault/server/internal/auth"
	"env-vault/server/internal/domain"
	"env-vault/server/internal/repository"
)

type AuthService interface {
	Register(ctx context.Context, email, password, name string) (*domain.User, error)
	Login(ctx context.Context, email, password, ip, userAgent string) (*domain.User, *domain.Session, string, string, error)
	Reauth(ctx context.Context, userID, password string) (string, error)
	Logout(ctx context.Context, sessionID string) error
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	GetCurrentUser(ctx context.Context, userID string) (*domain.User, error)
	ListSessions(ctx context.Context, userID string) ([]*domain.Session, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}

type authService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	tokenRepo   repository.RefreshTokenRepository
	resetRepo   repository.PasswordResetTokenRepository
	jwtService  *auth.JWTService
	argon2Cfg   auth.Argon2Config
	refreshTTL  time.Duration
}

func NewAuthService(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	tokenRepo repository.RefreshTokenRepository,
	resetRepo repository.PasswordResetTokenRepository,
	jwtService *auth.JWTService,
	argon2Cfg auth.Argon2Config,
	refreshTTL time.Duration,
) AuthService {
	return &authService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		tokenRepo:   tokenRepo,
		resetRepo:   resetRepo,
		jwtService:  jwtService,
		argon2Cfg:   argon2Cfg,
		refreshTTL:  refreshTTL,
	}
}

func (s *authService) Register(ctx context.Context, email, password, name string) (*domain.User, error) {
	if email == "" || password == "" || name == "" {
		return nil, domain.ErrInvalidInput
	}

	// Check if first user → make superadmin
	count, err := s.userRepo.CountUsers(ctx)
	if err != nil {
		return nil, err
	}

	hash, err := auth.HashPassword(password, s.argon2Cfg)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &domain.User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: hash,
		Name:         name,
		IsSuperAdmin: count == 0, // First user is superadmin
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password, ip, userAgent string) (*domain.User, *domain.Session, string, string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, nil, "", "", domain.ErrUnauthorized
		}
		return nil, nil, "", "", err
	}

	valid, err := auth.VerifyPassword(password, user.PasswordHash)
	if err != nil || !valid {
		return nil, nil, "", "", domain.ErrUnauthorized
	}

	// Create session
	now := time.Now()
	session := &domain.Session{
		ID:           uuid.NewString(),
		UserID:       user.ID,
		IPAddress:    ip,
		UserAgent:    userAgent,
		Device:       parseDevice(userAgent),
		Browser:      parseBrowser(userAgent),
		CreatedAt:    now,
		LastActivity: now,
	}
	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		return nil, nil, "", "", err
	}

	// Generate access token
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email, session.ID, user.IsSuperAdmin)
	if err != nil {
		return nil, nil, "", "", err
	}

	// Generate and store refresh token
	rawRefresh, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, nil, "", "", err
	}

	refreshTokenRecord := &domain.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		SessionID: session.ID,
		TokenHash: auth.HashToken(rawRefresh),
		ExpiresAt: now.Add(s.refreshTTL),
		CreatedAt: now,
	}
	if err := s.tokenRepo.CreateRefreshToken(ctx, refreshTokenRecord); err != nil {
		return nil, nil, "", "", err
	}

	return user, session, accessToken, rawRefresh, nil
}

func (s *authService) Reauth(ctx context.Context, userID, password string) (string, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return "", err
	}

	valid, err := auth.VerifyPassword(password, user.PasswordHash)
	if err != nil || !valid {
		return "", domain.ErrUnauthorized
	}

	sessionID, _ := ctx.Value(domain.CtxSessionID).(string)

	return s.jwtService.GenerateReauthToken(user.ID, user.Email, sessionID)
}

func (s *authService) Logout(ctx context.Context, sessionID string) error {
	if err := s.tokenRepo.RevokeAllForSession(ctx, sessionID); err != nil {
		return err
	}
	return s.sessionRepo.RevokeSession(ctx, sessionID)
}

func (s *authService) RefreshToken(ctx context.Context, rawToken string) (string, string, error) {
	hash := auth.HashToken(rawToken)
	tokenRecord, err := s.tokenRepo.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		return "", "", domain.ErrUnauthorized
	}

	if tokenRecord.RevokedAt != nil {
		return "", "", domain.ErrUnauthorized
	}
	if time.Now().After(tokenRecord.ExpiresAt) {
		return "", "", domain.ErrUnauthorized
	}

	// Check session is still active
	session, err := s.sessionRepo.GetSessionByID(ctx, tokenRecord.SessionID)
	if err != nil || session.RevokedAt != nil {
		return "", "", domain.ErrUnauthorized
	}

	// Get user
	user, err := s.userRepo.GetUserByID(ctx, tokenRecord.UserID)
	if err != nil {
		return "", "", domain.ErrUnauthorized
	}

	// Revoke old refresh token (rotation)
	if err := s.tokenRepo.RevokeRefreshToken(ctx, tokenRecord.ID); err != nil {
		return "", "", err
	}

	// Generate new access token
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email, session.ID, user.IsSuperAdmin)
	if err != nil {
		return "", "", err
	}

	// Generate new refresh token
	newRawRefresh, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}
	now := time.Now()
	newRefreshRecord := &domain.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		SessionID: session.ID,
		TokenHash: auth.HashToken(newRawRefresh),
		ExpiresAt: now.Add(s.refreshTTL),
		CreatedAt: now,
	}
	if err := s.tokenRepo.CreateRefreshToken(ctx, newRefreshRecord); err != nil {
		return "", "", err
	}

	// Update session activity
	_ = s.sessionRepo.UpdateLastActivity(ctx, session.ID)

	return accessToken, newRawRefresh, nil
}

func (s *authService) GetCurrentUser(ctx context.Context, userID string) (*domain.User, error) {
	return s.userRepo.GetUserByID(ctx, userID)
}

func (s *authService) ListSessions(ctx context.Context, userID string) ([]*domain.Session, error) {
	return s.sessionRepo.ListUserSessions(ctx, userID)
}

func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			// Don't reveal user existence
			return nil
		}
		return err
	}

	// Delete any existing tokens
	_ = s.resetRepo.DeleteTokensByUserID(ctx, user.ID)

	rawToken, err := auth.GenerateRefreshToken() // We can reuse the secure token generator
	if err != nil {
		return err
	}

	hash := auth.HashToken(rawToken)
	now := time.Now()
	token := &domain.PasswordResetToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		TokenHash: hash,
		ExpiresAt: now.Add(15 * time.Minute),
		CreatedAt: now,
	}

	if err := s.resetRepo.CreateToken(ctx, token); err != nil {
		return err
	}

	// In a real app, send email here. For now, log it.
	// log.Printf("MOCK EMAIL: Send password reset token '%s' to %s", rawToken, user.Email)
	return nil
}

func (s *authService) ResetPassword(ctx context.Context, rawToken, newPassword string) error {
	hash := auth.HashToken(rawToken)
	tokenRecord, err := s.resetRepo.GetTokenByHash(ctx, hash)
	if err != nil {
		return domain.ErrInvalidInput // Generic error so we don't leak token validity easily
	}

	if time.Now().After(tokenRecord.ExpiresAt) {
		return domain.ErrInvalidInput
	}

	user, err := s.userRepo.GetUserByID(ctx, tokenRecord.UserID)
	if err != nil {
		return domain.ErrInvalidInput
	}

	newHash, err := auth.HashPassword(newPassword, s.argon2Cfg)
	if err != nil {
		return err
	}

	// Update user password - we need an UpdateUser method. Wait, we don't have UpdateUser yet.
	// Oh, I need to add UpdateUser to userRepo.
	user.PasswordHash = newHash
	user.UpdatedAt = time.Now()
	
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return err
	}

	_ = s.resetRepo.DeleteTokensByUserID(ctx, user.ID)
	return nil
}

// Simple UA parsing helpers
func parseDevice(ua string) string {
	if ua == "" {
		return "unknown"
	}
	if containsAny(ua, "Mobile", "Android", "iPhone") {
		return "mobile"
	}
	if containsAny(ua, "Tablet", "iPad") {
		return "tablet"
	}
	return "desktop"
}

func parseBrowser(ua string) string {
	if ua == "" {
		return "unknown"
	}
	switch {
	case containsAny(ua, "Firefox"):
		return "Firefox"
	case containsAny(ua, "Edg"):
		return "Edge"
	case containsAny(ua, "Chrome"):
		return "Chrome"
	case containsAny(ua, "Safari"):
		return "Safari"
	default:
		return "other"
	}
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if len(s) >= len(sub) {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
		}
	}
	return false
}
