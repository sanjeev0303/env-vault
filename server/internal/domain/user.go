package domain

import "time"

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Name         string
	IsSuperAdmin bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type Session struct {
	ID           string
	UserID       string
	IPAddress    string
	UserAgent    string
	Device       string
	Browser      string
	CreatedAt    time.Time
	LastActivity time.Time
	RevokedAt    *time.Time
}

type RefreshToken struct {
	ID        string
	UserID    string
	SessionID string
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt *time.Time
}

type UserDevice struct {
	ID         string
	UserID     string
	DeviceName string
	DeviceType string
	Browser    string
	OS         string
	LastIP     string
	FirstSeen  time.Time
	LastSeen   time.Time
}

type PasswordResetToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}
