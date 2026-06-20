package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// InitDB initializes PostgreSQL database and creates tables if they don't exist.
func InitDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := createTables(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

type Table struct {
	Name       string
	Definition string
}

func createTables(db *sql.DB) error {
	tables := []Table{
		{
			Name: "projects",
			Definition: `(
				id TEXT PRIMARY KEY,
				name TEXT UNIQUE NOT NULL,
				created_at TIMESTAMP NOT NULL,
				updated_at TIMESTAMP NOT NULL
			)`,
		},
		{
			Name: "environments",
			Definition: `(
				id TEXT PRIMARY KEY,
				project_id TEXT NOT NULL,
				name TEXT NOT NULL,
				created_at TIMESTAMP NOT NULL,
				updated_at TIMESTAMP NOT NULL,
				FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
				UNIQUE(project_id, name)
			)`,
		},
		{
			Name: "secrets",
			Definition: `(
				id TEXT PRIMARY KEY,
				project_id TEXT NOT NULL,
				environment_id TEXT NOT NULL,
				key TEXT NOT NULL,
				value TEXT NOT NULL,
				created_at TIMESTAMP NOT NULL,
				updated_at TIMESTAMP NOT NULL,
				FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
				FOREIGN KEY (environment_id) REFERENCES environments(id) ON DELETE CASCADE,
				UNIQUE(environment_id, key)
			)`,
		},
	}

	for _, t := range tables {
		query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s %s;", t.Name, t.Definition)
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to create table %s: %w", t.Name, err)
		}
	}

	return nil
}
