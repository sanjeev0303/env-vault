package dto

import "time"

type CreateSecretReq struct {
	Key   string `json:"key" validate:"required,min=1,max=255"`
	Value string `json:"value" validate:"required"`
}

type UpdateSecretReq struct {
	Value string `json:"value" validate:"required"`
}

type RollbackReq struct {
	Version int `json:"version" validate:"required,min=1"`
}

// SecretMetadataResp - returned for list operations (no value).
type SecretMetadataResp struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"project_id"`
	EnvironmentID string    `json:"environment_id"`
	Key           string    `json:"key"`
	Version       int       `json:"version"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// SecretRevealResp - returned for single secret reveal (with value).
type SecretRevealResp struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"project_id"`
	EnvironmentID string    `json:"environment_id"`
	Key           string    `json:"key"`
	Value         string    `json:"value"`
	Version       int       `json:"version"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// SecretVersionResp - for version history.
type SecretVersionResp struct {
	ID        string    `json:"id"`
	SecretID  string    `json:"secret_id"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}

// Legacy compat (kept for .env format export)
type SecretResp struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"project_id"`
	EnvironmentID string    `json:"environment_id"`
	Key           string    `json:"key"`
	Value         string    `json:"value"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
