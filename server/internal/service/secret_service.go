package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"env-vault/server/internal/domain"
	"env-vault/server/internal/repository"
)

type secretService struct {
	secretRepo  repository.SecretRepository
	projectRepo repository.ProjectRepository
	encryptor   EncryptionService
}

func NewSecretService(
	secretRepo repository.SecretRepository,
	projectRepo repository.ProjectRepository,
	encryptor EncryptionService,
) SecretService {
	return &secretService{
		secretRepo:  secretRepo,
		projectRepo: projectRepo,
		encryptor:   encryptor,
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
	encryptedValue, err := s.encryptor.Encrypt(value)
	if err != nil {
		return nil, err
	}

	sec := &domain.Secret{
		ID:            uuid.NewString(),
		ProjectID:     projectID,
		EnvironmentID: envID,
		Key:           key,
		Value:         encryptedValue,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.secretRepo.CreateSecret(ctx, sec); err != nil {
		return nil, err
	}

	// Return decrypted secret to caller
	sec.Value = value
	return sec, nil
}

func (s *secretService) GetSecret(ctx context.Context, id string) (*domain.Secret, error) {
	sec, err := s.secretRepo.GetSecretByID(ctx, id)
	if err != nil {
		return nil, err
	}

	decryptedValue, err := s.encryptor.Decrypt(sec.Value)
	if err != nil {
		return nil, err
	}

	sec.Value = decryptedValue
	return sec, nil
}

func (s *secretService) ListSecrets(ctx context.Context, projectID, envID string) ([]*domain.Secret, error) {
	// Verify environment exists and matches projectID
	env, err := s.projectRepo.GetEnvironmentByID(ctx, envID)
	if err != nil {
		return nil, err
	}
	if env.ProjectID != projectID {
		return nil, domain.ErrInvalidInput
	}

	secrets, err := s.secretRepo.ListSecrets(ctx, envID)
	if err != nil {
		return nil, err
	}

	for _, sec := range secrets {
		decryptedValue, err := s.encryptor.Decrypt(sec.Value)
		if err != nil {
			return nil, err
		}
		sec.Value = decryptedValue
	}

	return secrets, nil
}

func (s *secretService) UpdateSecret(ctx context.Context, id, value string) (*domain.Secret, error) {
	sec, err := s.secretRepo.GetSecretByID(ctx, id)
	if err != nil {
		return nil, err
	}

	encryptedValue, err := s.encryptor.Encrypt(value)
	if err != nil {
		return nil, err
	}

	sec.Value = encryptedValue
	sec.UpdatedAt = time.Now()

	if err := s.secretRepo.UpdateSecret(ctx, sec); err != nil {
		return nil, err
	}

	sec.Value = value
	return sec, nil
}

func (s *secretService) DeleteSecret(ctx context.Context, id string) error {
	return s.secretRepo.DeleteSecret(ctx, id)
}

