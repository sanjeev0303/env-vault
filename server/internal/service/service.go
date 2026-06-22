package service

import (
	"context"

	"env-vault/server/internal/domain"
)

type SecretService interface {
	CreateSecret(ctx context.Context, projectID, envID, key, value string) (*domain.Secret, error)
	GetSecret(ctx context.Context, id string) (*domain.Secret, error)
	ListSecrets(ctx context.Context, projectID, envID string, cursor string, limit int) ([]*domain.Secret, string, error)
	ExportSecrets(ctx context.Context, projectID, envID string) (map[string]string, error)
	RevealSecret(ctx context.Context, id string) (*domain.Secret, error)
	UpdateSecret(ctx context.Context, id, value string) (*domain.Secret, error)
	DeleteSecret(ctx context.Context, id string) error
	GetSecretHistory(ctx context.Context, secretID string) ([]*domain.SecretVersion, error)
	RollbackSecret(ctx context.Context, secretID string, version int) (*domain.Secret, error)
}

type EncryptionService interface {
	Encrypt(plainText string) (*domain.Envelope, error)
	Decrypt(env *domain.Envelope) (string, error)
}
