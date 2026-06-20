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
		`CREATE TABLE projects (
			id TEXT PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);`,
		`CREATE TABLE environments (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			name TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
			UNIQUE(project_id, name)
		);`,
		`CREATE TABLE secrets (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			environment_id TEXT NOT NULL,
			key TEXT NOT NULL,
			value TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
			FOREIGN KEY (environment_id) REFERENCES environments(id) ON DELETE CASCADE,
			UNIQUE(environment_id, key)
		);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			db.Close()
			t.Fatalf("failed to execute migration query %q: %v", q, err)
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

	projectSvc := NewProjectService(projectRepo)
	secretSvc := NewSecretService(secretRepo, projectRepo, encryptor)

	ctx := context.Background()

	// 1. Create project
	proj, err := projectSvc.CreateProject(ctx, "Test Project")
	if err != nil {
		t.Fatalf("failed to create project: %v", err)
	}
	if proj.Name != "Test Project" {
		t.Errorf("expected project name %q, got %q", "Test Project", proj.Name)
	}

	// 2. Create duplicate project (should fail)
	_, err = projectSvc.CreateProject(ctx, "Test Project")
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
	if sec.Key != "DB_PASSWORD" || sec.Value != secretVal {
		t.Errorf("expected secret key/value %q/%q, got %q/%q", "DB_PASSWORD", secretVal, sec.Key, sec.Value)
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

	// 5. Get secret (should decrypt)
	fetchedSec, err := secretSvc.GetSecret(ctx, sec.ID)
	if err != nil {
		t.Fatalf("failed to get secret: %v", err)
	}
	if fetchedSec.Value != secretVal {
		t.Errorf("expected decrypted secret value %q, got %q", secretVal, fetchedSec.Value)
	}

	// 6. List secrets
	secList, err := secretSvc.ListSecrets(ctx, proj.ID, env.ID)
	if err != nil {
		t.Fatalf("failed to list secrets: %v", err)
	}
	if len(secList) != 1 || secList[0].Value != secretVal {
		t.Errorf("expected list with 1 decrypted secret, got %+v", secList)
	}

	// 7. Update secret
	newSecretVal := "updated-db-password"
	updatedSec, err := secretSvc.UpdateSecret(ctx, sec.ID, newSecretVal)
	if err != nil {
		t.Fatalf("failed to update secret: %v", err)
	}

	// Get secret again and verify update
	fetchedSecAfterUpdate, err := secretSvc.GetSecret(ctx, sec.ID)
	if err != nil {
		t.Fatalf("failed to get secret after update: %v", err)
	}
	if fetchedSecAfterUpdate.Value != newSecretVal {
		t.Errorf("expected updated value %q, got %q", newSecretVal, fetchedSecAfterUpdate.Value)
	}
	if updatedSec.Value != newSecretVal {
		t.Errorf("expected updated value in return object %q, got %q", newSecretVal, updatedSec.Value)
	}

	// 8. Delete secret
	err = secretSvc.DeleteSecret(ctx, sec.ID)
	if err != nil {
		t.Fatalf("failed to delete secret: %v", err)
	}

	_, err = secretSvc.GetSecret(ctx, sec.ID)
	if err == nil {
		t.Error("expected error getting deleted secret, got nil")
	}
}
