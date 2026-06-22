package service

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"

	"env-vault/server/internal/encryption"
	"env-vault/server/internal/repository/postgres"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		db.Close()
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	queries := []string{
		`CREATE TABLE users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			name TEXT NOT NULL,
			is_superadmin BOOLEAN NOT NULL DEFAULT FALSE,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			deleted_at DATETIME
		);`,
		`CREATE TABLE projects (
			id TEXT PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			owner_id TEXT,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			deleted_at DATETIME
		);`,
		`CREATE TABLE environments (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			name TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			deleted_at DATETIME,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
			UNIQUE(project_id, name)
		);`,
		`CREATE TABLE secrets (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			environment_id TEXT NOT NULL,
			key TEXT NOT NULL,
			value TEXT NOT NULL,
			version INTEGER NOT NULL DEFAULT 1,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			deleted_at DATETIME,
			encrypted_dek TEXT NOT NULL DEFAULT '',
			iv TEXT NOT NULL DEFAULT '',
			auth_tag TEXT NOT NULL DEFAULT '',
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
			FOREIGN KEY (environment_id) REFERENCES environments(id) ON DELETE CASCADE,
			UNIQUE(environment_id, key)
		);`,
		`CREATE TABLE secret_versions (
			id TEXT PRIMARY KEY,
			secret_id TEXT NOT NULL,
			version INTEGER NOT NULL,
			encrypted_value TEXT NOT NULL,
			encrypted_dek TEXT NOT NULL DEFAULT '',
			iv TEXT NOT NULL DEFAULT '',
			auth_tag TEXT NOT NULL DEFAULT '',
			created_by TEXT,
			created_at DATETIME NOT NULL,
			UNIQUE(secret_id, version),
			FOREIGN KEY (secret_id) REFERENCES secrets(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE project_members (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'viewer',
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			UNIQUE(project_id, user_id),
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE audit_logs (
			id TEXT PRIMARY KEY,
			project_id TEXT,
			user_id TEXT,
			action TEXT NOT NULL,
			target_id TEXT,
			details TEXT,
			ip_address TEXT,
			previous_hash TEXT NOT NULL,
			hash TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
		);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			db.Close()
			t.Fatalf("failed to execute migration query: %v", err)
		}
	}

	return db
}

func TestProjectAndSecretService(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	encryptor, err := encryption.NewAESEncryptionService("test-secret-key-123456789")
	if err != nil {
		t.Fatalf("failed to init encryptor: %v", err)
	}

	projectRepo := postgres.NewProjectRepository(db)
	secretRepo := postgres.NewSecretRepository(db)
	memberRepo := postgres.NewProjectMemberRepository(db)
	userRepo := postgres.NewUserRepository(db)
	auditRepo := postgres.NewAuditRepository(db)

	projectSvc := NewProjectService(projectRepo, memberRepo, userRepo)
	auditSvc := NewAuditService(auditRepo)
	secretSvc := NewSecretService(secretRepo, projectRepo, encryptor, auditSvc)

	ctx := context.Background()

	// 1. Create project (no owner for test)
	proj, err := projectSvc.CreateProject(ctx, "Test Project", "")
	if err != nil {
		t.Fatalf("failed to create project: %v", err)
	}
	if proj.Name != "Test Project" {
		t.Errorf("expected project name %q, got %q", "Test Project", proj.Name)
	}

	// 2. Create duplicate project (should fail)
	_, err = projectSvc.CreateProject(ctx, "Test Project", "")
	if err == nil {
		t.Error("expected error creating duplicate project, got nil")
	}

	// 3. Create environment
	env, err := projectSvc.CreateEnvironment(ctx, proj.ID, "production")
	if err != nil {
		t.Fatalf("failed to create environment: %v", err)
	}
	if env.Name != "production" {
		t.Errorf("expected env name %q, got %q", "production", env.Name)
	}

	// 4. Create secret
	secretVal := "secret-db-password"
	sec, err := secretSvc.CreateSecret(ctx, proj.ID, env.ID, "DB_PASSWORD", secretVal)
	if err != nil {
		t.Fatalf("failed to create secret: %v", err)
	}
	if sec.Key != "DB_PASSWORD" {
		t.Errorf("expected secret key %q, got %q", "DB_PASSWORD", sec.Key)
	}
	// Value should be empty in response (metadata only)
	if sec.Value != "" {
		t.Errorf("expected empty value in create response, got %q", sec.Value)
	}

	// Verify encryption at rest in DB
	var dbVal string
	err = db.QueryRow("SELECT value FROM secrets WHERE id = ?", sec.ID).Scan(&dbVal)
	if err != nil {
		t.Fatalf("failed to select secret value from database: %v", err)
	}
	if dbVal == secretVal {
		t.Errorf("value stored in database was not encrypted: got plain %q", dbVal)
	}

	// 5. List secrets (should return metadata only, no values)
	secList, _, err := secretSvc.ListSecrets(ctx, proj.ID, env.ID, "", 10)
	if err != nil {
		t.Fatalf("failed to list secrets: %v", err)
	}
	if len(secList) != 1 {
		t.Fatalf("expected 1 secret, got %d", len(secList))
	}
	if secList[0].Value != "" {
		t.Errorf("expected empty value in list response, got %q", secList[0].Value)
	}

	// 6. Reveal secret (should decrypt)
	revealed, err := secretSvc.RevealSecret(ctx, sec.ID)
	if err != nil {
		t.Fatalf("failed to reveal secret: %v", err)
	}
	if revealed.Value != secretVal {
		t.Errorf("expected revealed value %q, got %q", secretVal, revealed.Value)
	}

	// 7. Update secret (should create version history)
	newSecretVal := "updated-db-password"
	updatedSec, err := secretSvc.UpdateSecret(ctx, sec.ID, newSecretVal)
	if err != nil {
		t.Fatalf("failed to update secret: %v", err)
	}
	if updatedSec.Version != 2 {
		t.Errorf("expected version 2, got %d", updatedSec.Version)
	}

	// Verify updated value via reveal
	revealedAfterUpdate, err := secretSvc.RevealSecret(ctx, sec.ID)
	if err != nil {
		t.Fatalf("failed to reveal after update: %v", err)
	}
	if revealedAfterUpdate.Value != newSecretVal {
		t.Errorf("expected revealed value %q, got %q", newSecretVal, revealedAfterUpdate.Value)
	}

	// 8. Check version history
	versions, err := secretSvc.GetSecretHistory(ctx, sec.ID)
	if err != nil {
		t.Fatalf("failed to get secret history: %v", err)
	}
	if len(versions) != 1 {
		t.Errorf("expected 1 version in history, got %d", len(versions))
	}

	// 9. Rollback to version 1
	rolledBack, err := secretSvc.RollbackSecret(ctx, sec.ID, 1)
	if err != nil {
		t.Fatalf("failed to rollback secret: %v", err)
	}
	if rolledBack.Version != 3 {
		t.Errorf("expected version 3 after rollback, got %d", rolledBack.Version)
	}

	// Verify rollback restored original value
	revealedAfterRollback, err := secretSvc.RevealSecret(ctx, sec.ID)
	if err != nil {
		t.Fatalf("failed to reveal after rollback: %v", err)
	}
	if revealedAfterRollback.Value != secretVal {
		t.Errorf("expected rolled-back value %q, got %q", secretVal, revealedAfterRollback.Value)
	}

	// 10. Delete secret
	err = secretSvc.DeleteSecret(ctx, sec.ID)
	if err != nil {
		t.Fatalf("failed to delete secret: %v", err)
	}

	_, err = secretSvc.RevealSecret(ctx, sec.ID)
	if err == nil {
		t.Error("expected error revealing deleted secret, got nil")
	}
}
