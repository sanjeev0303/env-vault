package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"env-vault/server/internal/domain"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) CreateProject(ctx context.Context, p *domain.Project) error {
	query := `INSERT INTO projects (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.Name, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return domain.ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (r *ProjectRepository) GetProjectByID(ctx context.Context, id string) (*domain.Project, error) {
	query := `SELECT id, name, created_at, updated_at FROM projects WHERE id = $1`
	p := &domain.Project{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Name, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return p, nil
}

func (r *ProjectRepository) GetProjectByName(ctx context.Context, name string) (*domain.Project, error) {
	query := `SELECT id, name, created_at, updated_at FROM projects WHERE name = $1`
	p := &domain.Project{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(&p.ID, &p.Name, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return p, nil
}

func (r *ProjectRepository) ListProjects(ctx context.Context) ([]*domain.Project, error) {
	query := `SELECT id, name, created_at, updated_at FROM projects ORDER BY name ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*domain.Project
	for rows.Next() {
		p := &domain.Project{}
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *ProjectRepository) DeleteProject(ctx context.Context, id string) error {
	query := `DELETE FROM projects WHERE id = $1`
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

func (r *ProjectRepository) CreateEnvironment(ctx context.Context, env *domain.Environment) error {
	query := `INSERT INTO environments (id, project_id, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, env.ID, env.ProjectID, env.Name, env.CreatedAt, env.UpdatedAt)
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

func (r *ProjectRepository) GetEnvironmentByID(ctx context.Context, id string) (*domain.Environment, error) {
	query := `SELECT id, project_id, name, created_at, updated_at FROM environments WHERE id = $1`
	env := &domain.Environment{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&env.ID, &env.ProjectID, &env.Name, &env.CreatedAt, &env.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return env, nil
}

func (r *ProjectRepository) GetEnvironmentByName(ctx context.Context, projectID, name string) (*domain.Environment, error) {
	query := `SELECT id, project_id, name, created_at, updated_at FROM environments WHERE project_id = $1 AND name = $2`
	env := &domain.Environment{}
	err := r.db.QueryRowContext(ctx, query, projectID, name).Scan(&env.ID, &env.ProjectID, &env.Name, &env.CreatedAt, &env.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return env, nil
}

func (r *ProjectRepository) ListEnvironments(ctx context.Context, projectID string) ([]*domain.Environment, error) {
	query := `SELECT id, project_id, name, created_at, updated_at FROM environments WHERE project_id = $1 ORDER BY name ASC`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var envs []*domain.Environment
	for rows.Next() {
		env := &domain.Environment{}
		if err := rows.Scan(&env.ID, &env.ProjectID, &env.Name, &env.CreatedAt, &env.UpdatedAt); err != nil {
			return nil, err
		}
		envs = append(envs, env)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return envs, nil
}

func (r *ProjectRepository) DeleteEnvironment(ctx context.Context, id string) error {
	query := `DELETE FROM environments WHERE id = $1`
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
