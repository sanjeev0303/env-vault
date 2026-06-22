package repository

import (
	"context"

	"env-vault/server/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	CountUsers(ctx context.Context) (int, error)
}

type PasswordResetTokenRepository interface {
	CreateToken(ctx context.Context, token *domain.PasswordResetToken) error
	GetTokenByHash(ctx context.Context, hash string) (*domain.PasswordResetToken, error)
	DeleteTokensByUserID(ctx context.Context, userID string) error
}

type SessionRepository interface {
	CreateSession(ctx context.Context, session *domain.Session) error
	GetSessionByID(ctx context.Context, id string) (*domain.Session, error)
	UpdateLastActivity(ctx context.Context, id string) error
	RevokeSession(ctx context.Context, id string) error
	RevokeAllUserSessions(ctx context.Context, userID string) error
	ListUserSessions(ctx context.Context, userID string) ([]*domain.Session, error)
}

type RefreshTokenRepository interface {
	CreateRefreshToken(ctx context.Context, token *domain.RefreshToken) error
	GetRefreshTokenByHash(ctx context.Context, hash string) (*domain.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, id string) error
	RevokeAllForSession(ctx context.Context, sessionID string) error
}

type ProjectMemberRepository interface {
	AddMember(ctx context.Context, member *domain.ProjectMember) error
	GetMember(ctx context.Context, projectID, userID string) (*domain.ProjectMember, error)
	UpdateMemberRole(ctx context.Context, projectID, userID string, role domain.Role) error
	RemoveMember(ctx context.Context, projectID, userID string) error
	ListProjectMembers(ctx context.Context, projectID string) ([]*domain.ProjectMember, error)
	ListUserProjects(ctx context.Context, userID string) ([]string, error)
}

type AuditRepository interface {
	CreateLog(ctx context.Context, log *domain.AuditLog) error
	ListLogs(ctx context.Context, projectID string, cursor string, limit int) ([]*domain.AuditLog, string, error)
	GetLastHash(ctx context.Context) (string, error)
}
