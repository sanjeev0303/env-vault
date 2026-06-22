package postgres

import (
	"context"
	"database/sql"
	"errors"

	"env-vault/server/internal/domain"
)

type PasswordResetTokenRepository struct {
	db *sql.DB
}

func NewPasswordResetTokenRepository(db *sql.DB) *PasswordResetTokenRepository {
	return &PasswordResetTokenRepository{db: db}
}

func (r *PasswordResetTokenRepository) CreateToken(ctx context.Context, token *domain.PasswordResetToken) error {
	query := `INSERT INTO password_reset_tokens (id, user_id, token_hash, expires_at, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, token.ID, token.UserID, token.TokenHash, token.ExpiresAt, token.CreatedAt)
	return err
}

func (r *PasswordResetTokenRepository) GetTokenByHash(ctx context.Context, hash string) (*domain.PasswordResetToken, error) {
	query := `SELECT id, user_id, token_hash, expires_at, created_at FROM password_reset_tokens WHERE token_hash = $1`
	t := &domain.PasswordResetToken{}
	err := r.db.QueryRowContext(ctx, query, hash).Scan(&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return t, nil
}

func (r *PasswordResetTokenRepository) DeleteTokensByUserID(ctx context.Context, userID string) error {
	query := `DELETE FROM password_reset_tokens WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
