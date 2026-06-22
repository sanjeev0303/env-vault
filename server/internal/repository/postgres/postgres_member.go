package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"env-vault/server/internal/domain"
)

type ProjectMemberRepository struct {
	db *sql.DB
}

func NewProjectMemberRepository(db *sql.DB) *ProjectMemberRepository {
	return &ProjectMemberRepository{db: db}
}

func (r *ProjectMemberRepository) AddMember(ctx context.Context, m *domain.ProjectMember) error {
	query := `INSERT INTO project_members (id, project_id, user_id, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, m.ID, m.ProjectID, m.UserID, string(m.Role), m.CreatedAt, m.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			return domain.ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (r *ProjectMemberRepository) GetMember(ctx context.Context, projectID, userID string) (*domain.ProjectMember, error) {
	query := `SELECT id, project_id, user_id, role, created_at, updated_at FROM project_members WHERE project_id = $1 AND user_id = $2`
	m := &domain.ProjectMember{}
	var role string
	err := r.db.QueryRowContext(ctx, query, projectID, userID).Scan(&m.ID, &m.ProjectID, &m.UserID, &role, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	m.Role = domain.Role(role)
	return m, nil
}

func (r *ProjectMemberRepository) UpdateMemberRole(ctx context.Context, projectID, userID string, role domain.Role) error {
	query := `UPDATE project_members SET role = $1, updated_at = NOW() WHERE project_id = $2 AND user_id = $3`
	res, err := r.db.ExecContext(ctx, query, string(role), projectID, userID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ProjectMemberRepository) RemoveMember(ctx context.Context, projectID, userID string) error {
	query := `DELETE FROM project_members WHERE project_id = $1 AND user_id = $2`
	res, err := r.db.ExecContext(ctx, query, projectID, userID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ProjectMemberRepository) ListProjectMembers(ctx context.Context, projectID string) ([]*domain.ProjectMember, error) {
	query := `SELECT id, project_id, user_id, role, created_at, updated_at FROM project_members WHERE project_id = $1 ORDER BY created_at ASC`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*domain.ProjectMember
	for rows.Next() {
		m := &domain.ProjectMember{}
		var role string
		if err := rows.Scan(&m.ID, &m.ProjectID, &m.UserID, &role, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		m.Role = domain.Role(role)
		members = append(members, m)
	}
	return members, rows.Err()
}

func (r *ProjectMemberRepository) ListUserProjects(ctx context.Context, userID string) ([]string, error) {
	query := `SELECT project_id FROM project_members WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projectIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		projectIDs = append(projectIDs, id)
	}
	return projectIDs, rows.Err()
}
