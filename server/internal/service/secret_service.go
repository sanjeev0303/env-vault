package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"env-vault/server/internal/domain"
	"env-vault/server/internal/repository"
)

type secretService struct {
	secretRepo  repository.SecretRepository
	projectRepo repository.ProjectRepository
	encryptor   EncryptionService
	auditSvc    AuditService
}

func NewSecretService(
	secretRepo repository.SecretRepository,
	projectRepo repository.ProjectRepository,
	encryptor EncryptionService,
	auditSvc AuditService,
) SecretService {
	return &secretService{
		secretRepo:  secretRepo,
		projectRepo: projectRepo,
		encryptor:   encryptor,
		auditSvc:    auditSvc,
	}
}

func (s *secretService) CreateSecret(ctx context.Context, projectID, envID, key, value string) (*domain.Secret, error) {
	if projectID == "" || envID == "" || key == "" {
		return nil, domain.ErrInvalidInput
	}

	// Verify project and environment exist
	env, err := s.projectRepo.GetEnvironmentByID(ctx, envID)
	if err != nil {
		return nil, err
	}
	if env.ProjectID != projectID {
		return nil, domain.ErrInvalidInput
	}

	// Encrypt the value
	envelope, err := s.encryptor.Encrypt(value)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	sec := &domain.Secret{
		ID:            uuid.NewString(),
		ProjectID:     projectID,
		EnvironmentID: envID,
		Key:           key,
		Value:         envelope.EncryptedValue,
		EncryptedDEK:  envelope.EncryptedDEK,
		IV:            envelope.IV,
		AuthTag:       envelope.AuthTag,
		Version:       1,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.secretRepo.CreateSecret(ctx, sec); err != nil {
		return nil, err
	}

	_ = s.auditSvc.LogAction(ctx, domain.AuditCreate, projectID, sec.ID)

	// Return metadata only (no decrypted value)
	sec.Value = ""
	sec.EncryptedDEK = ""
	sec.IV = ""
	sec.AuthTag = ""
	return sec, nil
}

func (s *secretService) GetSecret(ctx context.Context, id string) (*domain.Secret, error) {
	sec, err := s.secretRepo.GetSecretByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Return metadata only - no decryption
	sec.Value = ""
	return sec, nil
}

// ListSecrets returns metadata only - never decrypts.
func (s *secretService) ListSecrets(ctx context.Context, projectID, envID string, cursor string, limit int) ([]*domain.Secret, string, error) {
	// Verify environment exists and matches projectID
	env, err := s.projectRepo.GetEnvironmentByID(ctx, envID)
	if err != nil {
		return nil, "", err
	}
	if env.ProjectID != projectID {
		return nil, "", domain.ErrInvalidInput
	}

	secrets, nextCursor, err := s.secretRepo.ListSecrets(ctx, envID, cursor, limit)
	if err != nil {
		return nil, "", err
	}

	// Strip encrypted values - return metadata only
	for _, sec := range secrets {
		sec.Value = ""
	}

	return secrets, nextCursor, nil
}

// RevealSecret decrypts a single secret value.
func (s *secretService) RevealSecret(ctx context.Context, id string) (*domain.Secret, error) {
	sec, err := s.secretRepo.GetSecretByID(ctx, id)
	if err != nil {
		return nil, err
	}

	env := &domain.Envelope{
		EncryptedValue: sec.Value,
		EncryptedDEK:   sec.EncryptedDEK,
		IV:             sec.IV,
		AuthTag:        sec.AuthTag,
	}

	decryptedValue, err := s.encryptor.Decrypt(env)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", domain.ErrInvalidInput)
	}

	_ = s.auditSvc.LogAction(ctx, domain.AuditReveal, sec.ProjectID, sec.ID)

	sec.Value = decryptedValue
	sec.EncryptedDEK = ""
	sec.IV = ""
	sec.AuthTag = ""
	return sec, nil
}

func (s *secretService) ExportSecrets(ctx context.Context, projectID, envID string) (map[string]string, error) {
	// Verify environment exists and matches projectID
	env, err := s.projectRepo.GetEnvironmentByID(ctx, envID)
	if err != nil {
		return nil, err
	}
	if env.ProjectID != projectID {
		return nil, domain.ErrInvalidInput
	}

	secrets, _, err := s.secretRepo.ListSecrets(ctx, envID, "", 10000)
	if err != nil {
		return nil, err
	}

	_ = s.auditSvc.LogAction(ctx, domain.AuditReveal, projectID, "EXPORT_ALL")

	result := make(map[string]string)
	for _, sec := range secrets {
		envelope := &domain.Envelope{
			EncryptedValue: sec.Value,
			EncryptedDEK:   sec.EncryptedDEK,
			IV:             sec.IV,
			AuthTag:        sec.AuthTag,
		}
		decryptedValue, err := s.encryptor.Decrypt(envelope)
		if err != nil {
			return nil, fmt.Errorf("decryption failed for secret %s: %w", sec.Key, err)
		}
		result[sec.Key] = decryptedValue
	}

	return result, nil
}

func (s *secretService) UpdateSecret(ctx context.Context, id, value string) (*domain.Secret, error) {
	sec, err := s.secretRepo.GetSecretByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Save version history before updating
	version := &domain.SecretVersion{
		ID:             uuid.NewString(),
		SecretID:       sec.ID,
		Version:        sec.Version,
		EncryptedValue: sec.Value,
		EncryptedDEK:   sec.EncryptedDEK,
		IV:             sec.IV,
		AuthTag:        sec.AuthTag,
		CreatedAt:      time.Time{}, // Handled by time.Now() below if needed or DB default
	}
	// Let's set CreatedAt properly
	version.CreatedAt = time.Now()

	if err := s.secretRepo.CreateSecretVersion(ctx, version); err != nil {
		return nil, err
	}

	envelope, err := s.encryptor.Encrypt(value)
	if err != nil {
		return nil, err
	}

	sec.Value = envelope.EncryptedValue
	sec.EncryptedDEK = envelope.EncryptedDEK
	sec.IV = envelope.IV
	sec.AuthTag = envelope.AuthTag
	sec.Version = sec.Version + 1
	sec.UpdatedAt = time.Now()

	if err := s.secretRepo.UpdateSecret(ctx, sec); err != nil {
		return nil, err
	}

	_ = s.auditSvc.LogAction(ctx, domain.AuditUpdate, sec.ProjectID, sec.ID)

	sec.Value = ""
	sec.EncryptedDEK = ""
	sec.IV = ""
	sec.AuthTag = ""
	return sec, nil
}

func (s *secretService) DeleteSecret(ctx context.Context, id string) error {
	sec, err := s.secretRepo.GetSecretByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.secretRepo.DeleteSecret(ctx, id)
	if err == nil {
		_ = s.auditSvc.LogAction(ctx, domain.AuditDelete, sec.ProjectID, sec.ID)
	}
	return err
}

func (s *secretService) GetSecretHistory(ctx context.Context, secretID string) ([]*domain.SecretVersion, error) {
	return s.secretRepo.ListSecretVersions(ctx, secretID)
}

func (s *secretService) RollbackSecret(ctx context.Context, secretID string, version int) (*domain.Secret, error) {
	// Get the version to restore
	v, err := s.secretRepo.GetSecretVersion(ctx, secretID, version)
	if err != nil {
		return nil, err
	}

	// Get current secret
	sec, err := s.secretRepo.GetSecretByID(ctx, secretID)
	if err != nil {
		return nil, err
	}

	// Save current as a version before rollback
	currentVersion := &domain.SecretVersion{
		ID:             uuid.NewString(),
		SecretID:       sec.ID,
		Version:        sec.Version,
		EncryptedValue: sec.Value,
		EncryptedDEK:   sec.EncryptedDEK,
		IV:             sec.IV,
		AuthTag:        sec.AuthTag,
		CreatedAt:      time.Now(),
	}
	if err := s.secretRepo.CreateSecretVersion(ctx, currentVersion); err != nil {
		return nil, err
	}

	// Restore the old version's encrypted value and envelope fields
	sec.Value = v.EncryptedValue
	sec.EncryptedDEK = v.EncryptedDEK
	sec.IV = v.IV
	sec.AuthTag = v.AuthTag
	sec.Version = sec.Version + 1
	sec.UpdatedAt = time.Now()

	if err := s.secretRepo.UpdateSecret(ctx, sec); err != nil {
		return nil, err
	}

	_ = s.auditSvc.LogAction(ctx, domain.AuditRollback, sec.ProjectID, sec.ID)

	sec.Value = ""
	sec.EncryptedDEK = ""
	sec.IV = ""
	sec.AuthTag = ""
	return sec, nil
}
