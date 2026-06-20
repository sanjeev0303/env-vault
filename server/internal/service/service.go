package service

import (
	"context"

	"env-vault/server/internal/domain"
)

type ProjectService interface {
	CreateProject(ctx context.Context, name string) (*domain.Project, error)
	GetProject(ctx context.Context, id string) (*domain.Project, error)
	GetProjectByName(ctx context.Context, name string) (*domain.Project, error)
	ListProjects(ctx context.Context) ([]*domain.Project, error)
	DeleteProject(ctx context.Context, id string) error

	CreateEnvironment(ctx context.Context, projectID, name string) (*domain.Environment, error)
	GetEnvironment(ctx context.Context, id string) (*domain.Environment, error)
	GetEnvironmentByName(ctx context.Context, projectID, name string) (*domain.Environment, error)
	ListEnvironments(ctx context.Context, projectID string) ([]*domain.Environment, error)
	DeleteEnvironment(ctx context.Context, id string) error
}

type SecretService interface {
	CreateSecret(ctx context.Context, projectID, envID, key, value string) (*domain.Secret, error)
	GetSecret(ctx context.Context, id string) (*domain.Secret, error)
	ListSecrets(ctx context.Context, projectID, envID string) ([]*domain.Secret, error)
	UpdateSecret(ctx context.Context, id, value string) (*domain.Secret, error)
	DeleteSecret(ctx context.Context, id string) error
}

type EncryptionService interface {
	Encrypt(plainText string) (string, error)
	Decrypt(cipherText string) (string, error)
}
