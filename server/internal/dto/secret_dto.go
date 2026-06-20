package dto

import "time"

type CreateSecretReq struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UpdateSecretReq struct {
	Value string `json:"value"`
}

type SecretResp struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"project_id"`
	EnvironmentID string    `json:"environment_id"`
	Key           string    `json:"key"`
	Value         string    `json:"value"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
