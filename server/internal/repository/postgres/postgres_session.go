package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"env-vault/server/internal/domain"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateSession(ctx context.Context, s *domain.Session) error {
	query := `INSERT INTO sessions (id, user_id, ip_address, user_agent, device, browser, created_at, last_activity) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, s.ID, s.UserID, s.IPAddress, s.UserAgent, s.Device, s.Browser, s.CreatedAt, s.LastActivity)
	return err
}

func (r *SessionRepository) GetSessionByID(ctx context.Context, id string) (*domain.Session, error) {
	query := `SELECT id, user_id, ip_address, user_agent, device, browser, created_at, last_activity, revoked_at FROM sessions WHERE id = $1`
	s := &domain.Session{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.UserID, &s.IPAddress, &s.UserAgent, &s.Device, &s.Browser, &s.CreatedAt, &s.LastActivity, &s.RevokedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return s, nil
}

func (r *SessionRepository) UpdateLastActivity(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE sessions SET last_activity = $1 WHERE id = $2", time.Now(), id)
	return err
}

func (r *SessionRepository) RevokeSession(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE sessions SET revoked_at = $1 WHERE id = $2", time.Now(), id)
	return err
}

func (r *SessionRepository) RevokeAllUserSessions(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE sessions SET revoked_at = $1 WHERE user_id = $2 AND revoked_at IS NULL", time.Now(), userID)
	return err
}

func (r *SessionRepository) ListUserSessions(ctx context.Context, userID string) ([]*domain.Session, error) {
	query := `SELECT id, user_id, ip_address, user_agent, device, browser, created_at, last_activity, revoked_at FROM sessions WHERE user_id = $1 ORDER BY last_activity DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*domain.Session
	for rows.Next() {
		s := &domain.Session{}
		if err := rows.Scan(&s.ID, &s.UserID, &s.IPAddress, &s.UserAgent, &s.Device, &s.Browser, &s.CreatedAt, &s.LastActivity, &s.RevokedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}

type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) CreateRefreshToken(ctx context.Context, t *domain.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (id, user_id, session_id, token_hash, expires_at, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, t.ID, t.UserID, t.SessionID, t.TokenHash, t.ExpiresAt, t.CreatedAt)
	return err
}

func (r *RefreshTokenRepository) GetRefreshTokenByHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	query := `SELECT id, user_id, session_id, token_hash, expires_at, created_at, revoked_at FROM refresh_tokens WHERE token_hash = $1`
	t := &domain.RefreshToken{}
	err := r.db.QueryRowContext(ctx, query, hash).Scan(&t.ID, &t.UserID, &t.SessionID, &t.TokenHash, &t.ExpiresAt, &t.CreatedAt, &t.RevokedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return t, nil
}

func (r *RefreshTokenRepository) RevokeRefreshToken(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE refresh_tokens SET revoked_at = $1 WHERE id = $2", time.Now(), id)
	return err
}

func (r *RefreshTokenRepository) RevokeAllForSession(ctx context.Context, sessionID string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE refresh_tokens SET revoked_at = $1 WHERE session_id = $2 AND revoked_at IS NULL", time.Now(), sessionID)
	return err
}
