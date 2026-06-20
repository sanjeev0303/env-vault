package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"env-vault/server/internal/domain"
	"env-vault/server/internal/repository"
)

type projectService struct {
	repo repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) ProjectService {
	return &projectService{repo: repo}
}

func (s *projectService) CreateProject(ctx context.Context, name string) (*domain.Project, error) {
	if name == "" {
		return nil, domain.ErrInvalidInput
	}

	p := &domain.Project{
		ID:        uuid.NewString(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateProject(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *projectService) GetProject(ctx context.Context, id string) (*domain.Project, error) {
	if id == "" {
		return nil, domain.ErrInvalidInput
	}
	return s.repo.GetProjectByID(ctx, id)
}

func (s *projectService) GetProjectByName(ctx context.Context, name string) (*domain.Project, error) {
	if name == "" {
		return nil, domain.ErrInvalidInput
	}
	return s.repo.GetProjectByName(ctx, name)
}

func (s *projectService) ListProjects(ctx context.Context) ([]*domain.Project, error) {
	return s.repo.ListProjects(ctx)
}

func (s *projectService) DeleteProject(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalidInput
	}
	return s.repo.DeleteProject(ctx, id)
}

func (s *projectService) CreateEnvironment(ctx context.Context, projectID, name string) (*domain.Environment, error) {
	if projectID == "" || name == "" {
		return nil, domain.ErrInvalidInput
	}

	// Verify project exists
	_, err := s.repo.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	env := &domain.Environment{
		ID:        uuid.NewString(),
		ProjectID: projectID,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateEnvironment(ctx, env); err != nil {
		return nil, err
	}

	return env, nil
}

func (s *projectService) GetEnvironment(ctx context.Context, id string) (*domain.Environment, error) {
	if id == "" {
		return nil, domain.ErrInvalidInput
	}
	return s.repo.GetEnvironmentByID(ctx, id)
}

func (s *projectService) GetEnvironmentByName(ctx context.Context, projectID, name string) (*domain.Environment, error) {
	if projectID == "" || name == "" {
		return nil, domain.ErrInvalidInput
	}
	return s.repo.GetEnvironmentByName(ctx, projectID, name)
}

func (s *projectService) ListEnvironments(ctx context.Context, projectID string) ([]*domain.Environment, error) {
	if projectID == "" {
		return nil, domain.ErrInvalidInput
	}
	return s.repo.ListEnvironments(ctx, projectID)
}

func (s *projectService) DeleteEnvironment(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalidInput
	}
	return s.repo.DeleteEnvironment(ctx, id)
}
