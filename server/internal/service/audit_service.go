package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"env-vault/server/internal/domain"
	"env-vault/server/internal/repository"
)

type AuditService interface {
	LogAction(ctx context.Context, action domain.AuditAction, projectID, secretID string) error
}

type auditService struct {
	auditRepo repository.AuditRepository
}

func NewAuditService(auditRepo repository.AuditRepository) AuditService {
	return &auditService{
		auditRepo: auditRepo,
	}
}

func (s *auditService) LogAction(ctx context.Context, action domain.AuditAction, projectID, secretID string) error {
	userID, _ := ctx.Value(domain.CtxUserID).(string)
	if userID == "" {
		userID = "system"
	}
	ip, _ := ctx.Value(domain.CtxIPAddress).(string)
	ua, _ := ctx.Value(domain.CtxUserAgent).(string)

	log := &domain.AuditLog{
		ID:        uuid.NewString(),
		UserID:    userID,
		ProjectID: projectID,
		SecretID:  secretID,
		Action:    action,
		IPAddress: ip,
		UserAgent: ua,
		CreatedAt: time.Now(),
	}

	return s.auditRepo.CreateLog(ctx, log)
}
