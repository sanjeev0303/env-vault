package domain

import (
	"time"
)

type Secret struct {
	ID            string
	ProjectID     string
	EnvironmentID string
	Key           string
	Value         string // plain text in services, cipher text in repositories
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
