package domain

import "time"

type AuditAction string

const (
	AuditReveal   AuditAction = "REVEAL"
	AuditCreate   AuditAction = "CREATE"
	AuditUpdate   AuditAction = "UPDATE"
	AuditDelete   AuditAction = "DELETE"
	AuditRollback AuditAction = "ROLLBACK"
	AuditExport   AuditAction = "EXPORT"
)

type AuditLog struct {
	ID           string
	UserID       string
	ProjectID    string
	SecretID     string
	Action       AuditAction
	IPAddress    string
	UserAgent    string
	PreviousHash string
	CurrentHash  string
	CreatedAt    time.Time
}

type SecretVersion struct {
	ID             string
	SecretID       string
	Version        int
	EncryptedValue string
	EncryptedDEK   string
	IV             string
	AuthTag        string
	CreatedBy      string
	CreatedAt      time.Time
}
