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
	query := `INSERT INTO secrets (id, project_id, environment_id, key, value, encrypted_dek, iv, auth_tag, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.ExecContext(ctx, query, s.ID, s.ProjectID, s.EnvironmentID, s.Key, s.Value, s.EncryptedDEK, s.IV, s.AuthTag, s.Version, s.CreatedAt, s.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "duplicate") {
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
	query := `SELECT id, project_id, environment_id, key, value, encrypted_dek, iv, auth_tag, version, created_at, updated_at FROM secrets WHERE id = $1`
	s := &domain.Secret{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.ProjectID, &s.EnvironmentID, &s.Key, &s.Value, &s.EncryptedDEK, &s.IV, &s.AuthTag, &s.Version, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return s, nil
}

func (r *SecretRepository) GetSecretByKey(ctx context.Context, envID string, key string) (*domain.Secret, error) {
	query := `SELECT id, project_id, environment_id, key, value, encrypted_dek, iv, auth_tag, version, created_at, updated_at FROM secrets WHERE environment_id = $1 AND key = $2`
	s := &domain.Secret{}
	err := r.db.QueryRowContext(ctx, query, envID, key).Scan(&s.ID, &s.ProjectID, &s.EnvironmentID, &s.Key, &s.Value, &s.EncryptedDEK, &s.IV, &s.AuthTag, &s.Version, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return s, nil
}

func (r *SecretRepository) ListSecrets(ctx context.Context, envID string, cursor string, limit int) ([]*domain.Secret, string, error) {
	var query string
	var args []interface{}
	
	if cursor == "" {
		query = `SELECT id, project_id, environment_id, key, value, encrypted_dek, iv, auth_tag, version, created_at, updated_at FROM secrets WHERE environment_id = $1 ORDER BY key ASC LIMIT $2`
		args = []interface{}{envID, limit + 1}
	} else {
		query = `SELECT id, project_id, environment_id, key, value, encrypted_dek, iv, auth_tag, version, created_at, updated_at FROM secrets WHERE environment_id = $1 AND key > $2 ORDER BY key ASC LIMIT $3`
		args = []interface{}{envID, cursor, limit + 1}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var secrets []*domain.Secret
	for rows.Next() {
		s := &domain.Secret{}
		if err := rows.Scan(&s.ID, &s.ProjectID, &s.EnvironmentID, &s.Key, &s.Value, &s.EncryptedDEK, &s.IV, &s.AuthTag, &s.Version, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, "", err
		}
		secrets = append(secrets, s)
	}
	if err := rows.Err(); err != nil {
		return nil, "", err
	}

	nextCursor := ""
	if len(secrets) > limit {
		nextCursor = secrets[limit].Key
		secrets = secrets[:limit]
	}

	return secrets, nextCursor, nil
}

func (r *SecretRepository) UpdateSecret(ctx context.Context, s *domain.Secret) error {
	query := `UPDATE secrets SET value = $1, encrypted_dek = $2, iv = $3, auth_tag = $4, version = $5, updated_at = $6 WHERE id = $7`
	res, err := r.db.ExecContext(ctx, query, s.Value, s.EncryptedDEK, s.IV, s.AuthTag, s.Version, s.UpdatedAt, s.ID)
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

// Versioning
func (r *SecretRepository) CreateSecretVersion(ctx context.Context, v *domain.SecretVersion) error {
	query := `INSERT INTO secret_versions (id, secret_id, version, encrypted_value, encrypted_dek, iv, auth_tag, created_by, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query, v.ID, v.SecretID, v.Version, v.EncryptedValue, v.EncryptedDEK, v.IV, v.AuthTag, v.CreatedBy, v.CreatedAt)
	return err
}

func (r *SecretRepository) ListSecretVersions(ctx context.Context, secretID string) ([]*domain.SecretVersion, error) {
	query := `SELECT id, secret_id, version, created_at FROM secret_versions WHERE secret_id = $1 ORDER BY version DESC`
	rows, err := r.db.QueryContext(ctx, query, secretID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []*domain.SecretVersion
	for rows.Next() {
		v := &domain.SecretVersion{}
		if err := rows.Scan(&v.ID, &v.SecretID, &v.Version, &v.CreatedAt); err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}
	return versions, rows.Err()
}

func (r *SecretRepository) GetSecretVersion(ctx context.Context, secretID string, version int) (*domain.SecretVersion, error) {
	query := `SELECT id, secret_id, version, encrypted_value, encrypted_dek, iv, auth_tag, created_by, created_at FROM secret_versions WHERE secret_id = $1 AND version = $2`
	v := &domain.SecretVersion{}
	err := r.db.QueryRowContext(ctx, query, secretID, version).Scan(&v.ID, &v.SecretID, &v.Version, &v.EncryptedValue, &v.EncryptedDEK, &v.IV, &v.AuthTag, &v.CreatedBy, &v.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return v, nil
}
