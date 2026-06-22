package database

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
)

type Migration struct {
	Version int
	Name    string
	SQL     string
}

var migrations = []Migration{
	{
		Version: 1,
		Name:    "create_users_table",
		SQL: `
CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY,
	email TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	name TEXT NOT NULL,
	is_superadmin BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
`,
	},
	{
		Version: 2,
		Name:    "create_sessions_table",
		SQL: `
CREATE TABLE IF NOT EXISTS sessions (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	ip_address TEXT NOT NULL DEFAULT '',
	user_agent TEXT NOT NULL DEFAULT '',
	device TEXT NOT NULL DEFAULT '',
	browser TEXT NOT NULL DEFAULT '',
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	last_activity TIMESTAMP NOT NULL DEFAULT NOW(),
	revoked_at TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
`,
	},
	{
		Version: 3,
		Name:    "create_refresh_tokens_table",
		SQL: `
CREATE TABLE IF NOT EXISTS refresh_tokens (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	session_id TEXT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
	token_hash TEXT NOT NULL,
	expires_at TIMESTAMP NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	revoked_at TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_session_id ON refresh_tokens(session_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
`,
	},
	{
		Version: 4,
		Name:    "create_user_devices_table",
		SQL: `
CREATE TABLE IF NOT EXISTS user_devices (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	device_name TEXT NOT NULL DEFAULT '',
	device_type TEXT NOT NULL DEFAULT '',
	browser TEXT NOT NULL DEFAULT '',
	os TEXT NOT NULL DEFAULT '',
	last_ip TEXT NOT NULL DEFAULT '',
	first_seen TIMESTAMP NOT NULL DEFAULT NOW(),
	last_seen TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_user_devices_user_id ON user_devices(user_id);
`,
	},
	{
		Version: 5,
		Name:    "create_project_members_table",
		SQL: `
CREATE TABLE IF NOT EXISTS project_members (
	id TEXT PRIMARY KEY,
	project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	role TEXT NOT NULL DEFAULT 'viewer',
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	UNIQUE(project_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_project_members_user_id ON project_members(user_id);
CREATE INDEX IF NOT EXISTS idx_project_members_project_id ON project_members(project_id);
`,
	},
	{
		Version: 6,
		Name:    "add_owner_to_projects",
		SQL: `
DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='projects' AND column_name='owner_id') THEN
		ALTER TABLE projects ADD COLUMN owner_id TEXT;
	END IF;
	IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='projects' AND column_name='deleted_at') THEN
		ALTER TABLE projects ADD COLUMN deleted_at TIMESTAMP;
	END IF;
END $$;
`,
	},
	{
		Version: 7,
		Name:    "add_soft_delete_to_envs_and_secrets",
		SQL: `
DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='environments' AND column_name='deleted_at') THEN
		ALTER TABLE environments ADD COLUMN deleted_at TIMESTAMP;
	END IF;
	IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='secrets' AND column_name='deleted_at') THEN
		ALTER TABLE secrets ADD COLUMN deleted_at TIMESTAMP;
	END IF;
END $$;
`,
	},
	{
		Version: 8,
		Name:    "add_envelope_encryption_columns",
		SQL: `
DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='secrets' AND column_name='encrypted_dek') THEN
		ALTER TABLE secrets ADD COLUMN encrypted_dek TEXT NOT NULL DEFAULT '';
	END IF;
	IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='secrets' AND column_name='iv') THEN
		ALTER TABLE secrets ADD COLUMN iv TEXT NOT NULL DEFAULT '';
	END IF;
	IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='secrets' AND column_name='auth_tag') THEN
		ALTER TABLE secrets ADD COLUMN auth_tag TEXT NOT NULL DEFAULT '';
	END IF;
	IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='secrets' AND column_name='version') THEN
		ALTER TABLE secrets ADD COLUMN version INTEGER NOT NULL DEFAULT 1;
	END IF;
END $$;
`,
	},
	{
		Version: 9,
		Name:    "create_secret_versions_table",
		SQL: `
CREATE TABLE IF NOT EXISTS secret_versions (
	id TEXT PRIMARY KEY,
	secret_id TEXT NOT NULL REFERENCES secrets(id) ON DELETE CASCADE,
	version INTEGER NOT NULL,
	encrypted_value TEXT NOT NULL,
	encrypted_dek TEXT NOT NULL,
	iv TEXT NOT NULL,
	auth_tag TEXT NOT NULL,
	created_by TEXT,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	UNIQUE(secret_id, version)
);
CREATE INDEX IF NOT EXISTS idx_secret_versions_secret_id ON secret_versions(secret_id);
`,
	},
	{
		Version: 10,
		Name:    "create_secret_access_logs_table",
		SQL: `
CREATE TABLE IF NOT EXISTS secret_access_logs (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL,
	project_id TEXT NOT NULL,
	secret_id TEXT,
	action TEXT NOT NULL,
	ip_address TEXT NOT NULL DEFAULT '',
	user_agent TEXT NOT NULL DEFAULT '',
	previous_hash TEXT NOT NULL DEFAULT '',
	current_hash TEXT NOT NULL DEFAULT '',
	created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON secret_access_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_project_id ON secret_access_logs(project_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON secret_access_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON secret_access_logs(action);
`,
	},
	{
		Version: 11,
		Name:    "add_indexes_to_existing_tables",
		SQL: `
CREATE INDEX IF NOT EXISTS idx_secrets_env_id ON secrets(environment_id);
CREATE INDEX IF NOT EXISTS idx_secrets_project_id ON secrets(project_id);
CREATE INDEX IF NOT EXISTS idx_environments_project_id ON environments(project_id);
`,
	},
	{
		Version: 12,
		Name:    "create_password_reset_tokens",
		SQL: `
CREATE TABLE IF NOT EXISTS password_reset_tokens (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	token_hash TEXT NOT NULL,
	expires_at TIMESTAMP NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_token_hash ON password_reset_tokens(token_hash);
`,
	},
}

// RunMigrations executes all pending migrations in order.
func RunMigrations(db *sql.DB) error {
	if err := ensureMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	applied, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	for _, m := range migrations {
		if applied[m.Version] {
			continue
		}

		log.Printf("Running migration %d: %s", m.Version, m.Name)

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %d: %w", m.Version, err)
		}

		if _, err := tx.Exec(m.SQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %d (%s): %w", m.Version, m.Name, err)
		}

		if _, err := tx.Exec(
			"INSERT INTO schema_migrations (version, name) VALUES ($1, $2)",
			m.Version, m.Name,
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %d: %w", m.Version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", m.Version, err)
		}

		log.Printf("Migration %d applied successfully", m.Version)
	}

	return nil
}

func ensureMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

func getAppliedMigrations(db *sql.DB) (map[int]bool, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		applied[v] = true
	}
	return applied, rows.Err()
}
