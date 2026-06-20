package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"env-vault/server/internal/domain"
)

type SecretRepository struct {
	db *sql.DB
}

func NewSecretRepository(db *sql.DB) *SecretRepository {
	return &SecretRepository{db: db}
}

func (r *SecretRepository) CreateSecret(ctx context.Context, s *domain.Secret) error {
	query := `INSERT INTO secrets (id, project_id, environment_id, key, value, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query, s.ID, s.ProjectID, s.EnvironmentID, s.Key, s.Value, s.CreatedAt, s.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return domain.ErrAlreadyExists
		}
		if strings.Contains(err.Error(), "foreign key") || strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			return domain.ErrNotFound
		}
		return err
	}
	return nil
}

func (r *SecretRepository) GetSecretByID(ctx context.Context, id string) (*domain.Secret, error) {
	query := `SELECT id, project_id, environment_id, key, value, created_at, updated_at FROM secrets WHERE id = $1`
	s := &domain.Secret{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.ProjectID, &s.EnvironmentID, &s.Key, &s.Value, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return s, nil
}

func (r *SecretRepository) GetSecretByKey(ctx context.Context, envID string, key string) (*domain.Secret, error) {
	query := `SELECT id, project_id, environment_id, key, value, created_at, updated_at FROM secrets WHERE environment_id = $1 AND key = $2`
	s := &domain.Secret{}
	err := r.db.QueryRowContext(ctx, query, envID, key).Scan(&s.ID, &s.ProjectID, &s.EnvironmentID, &s.Key, &s.Value, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return s, nil
}

func (r *SecretRepository) ListSecrets(ctx context.Context, envID string) ([]*domain.Secret, error) {
	query := `SELECT id, project_id, environment_id, key, value, created_at, updated_at FROM secrets WHERE environment_id = $1 ORDER BY key ASC`
	rows, err := r.db.QueryContext(ctx, query, envID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var secrets []*domain.Secret
	for rows.Next() {
		s := &domain.Secret{}
		if err := rows.Scan(&s.ID, &s.ProjectID, &s.EnvironmentID, &s.Key, &s.Value, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		secrets = append(secrets, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return secrets, nil
}

func (r *SecretRepository) UpdateSecret(ctx context.Context, s *domain.Secret) error {
	query := `UPDATE secrets SET value = $1, updated_at = $2 WHERE id = $3`
	res, err := r.db.ExecContext(ctx, query, s.Value, s.UpdatedAt, s.ID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *SecretRepository) DeleteSecret(ctx context.Context, id string) error {
	query := `DELETE FROM secrets WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}
