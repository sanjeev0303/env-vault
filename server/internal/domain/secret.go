package domain

import (
	"time"
)

type Envelope struct {
	EncryptedValue string
	EncryptedDEK   string
	IV             string
	AuthTag        string
}

type Secret struct {
	ID            string
	ProjectID     string
	EnvironmentID string
	Key           string
	Value         string // encrypted in DB, plain only during reveal
	EncryptedDEK  string
	IV            string
	AuthTag       string
	Version       int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
