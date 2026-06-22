package repository

import (
	"context"

	"env-vault/server/internal/domain"
)

type ProjectRepository interface {
	CreateProject(ctx context.Context, project *domain.Project) error
	GetProjectByID(ctx context.Context, id string) (*domain.Project, error)
	GetProjectByName(ctx context.Context, name string) (*domain.Project, error)
	ListProjects(ctx context.Context) ([]*domain.Project, error)
	DeleteProject(ctx context.Context, id string) error

	CreateEnvironment(ctx context.Context, env *domain.Environment) error
	GetEnvironmentByID(ctx context.Context, id string) (*domain.Environment, error)
	GetEnvironmentByName(ctx context.Context, projectID, name string) (*domain.Environment, error)
	ListEnvironments(ctx context.Context, projectID string) ([]*domain.Environment, error)
	DeleteEnvironment(ctx context.Context, id string) error
}

type SecretRepository interface {
	CreateSecret(ctx context.Context, secret *domain.Secret) error
	GetSecretByID(ctx context.Context, id string) (*domain.Secret, error)
	GetSecretByKey(ctx context.Context, envID string, key string) (*domain.Secret, error)
	ListSecrets(ctx context.Context, envID string, cursor string, limit int) ([]*domain.Secret, string, error)
	UpdateSecret(ctx context.Context, secret *domain.Secret) error
	DeleteSecret(ctx context.Context, id string) error

	// Versioning
	CreateSecretVersion(ctx context.Context, version *domain.SecretVersion) error
	ListSecretVersions(ctx context.Context, secretID string) ([]*domain.SecretVersion, error)
	GetSecretVersion(ctx context.Context, secretID string, version int) (*domain.SecretVersion, error)
}
