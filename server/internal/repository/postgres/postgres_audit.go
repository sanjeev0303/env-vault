package postgres

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"env-vault/server/internal/domain"
)

type AuditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) CreateLog(ctx context.Context, log *domain.AuditLog) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var previousHash string
	err = tx.QueryRowContext(ctx, "SELECT current_hash FROM secret_access_logs ORDER BY created_at DESC LIMIT 1 FOR UPDATE").Scan(&previousHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			previousHash = "0000000000000000000000000000000000000000000000000000000000000000"
		} else {
			return err
		}
	}

	log.PreviousHash = previousHash
	
	// Format: "action|projectID|secretID|userID|prevHash"
	data := fmt.Sprintf("%s|%s|%s|%s|%s", log.Action, log.ProjectID, log.SecretID, log.UserID, log.PreviousHash)
	hash := sha256.Sum256([]byte(data))
	log.CurrentHash = hex.EncodeToString(hash[:])

	query := `INSERT INTO secret_access_logs 
		(id, user_id, project_id, secret_id, action, ip_address, user_agent, previous_hash, current_hash, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	
	_, err = tx.ExecContext(ctx, query,
		log.ID, log.UserID, log.ProjectID, log.SecretID, log.Action,
		log.IPAddress, log.UserAgent, log.PreviousHash, log.CurrentHash, log.CreatedAt,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *AuditRepository) GetLastHash(ctx context.Context) (string, error) {
	var hash string
	err := r.db.QueryRowContext(ctx, "SELECT current_hash FROM secret_access_logs ORDER BY created_at DESC LIMIT 1").Scan(&hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "0000000000000000000000000000000000000000000000000000000000000000", nil
		}
		return "", err
	}
	return hash, nil
}

func (r *AuditRepository) ListLogs(ctx context.Context, projectID string, cursor string, limit int) ([]*domain.AuditLog, string, error) {
	var query string
	var args []interface{}

	if cursor == "" {
		query = `SELECT id, user_id, project_id, secret_id, action, ip_address, user_agent, previous_hash, current_hash, created_at FROM secret_access_logs WHERE project_id = $1 ORDER BY created_at DESC LIMIT $2`
		args = []interface{}{projectID, limit + 1}
	} else {
		// Use ID as cursor (since created_at might not be unique if timestamps are identical, wait, ID is UUID.
		// For true cursor, usually (created_at, id) < (cursor_created_at, cursor_id).
		// We'll use ID and fetch created_at inside the query or just use ID, but ordering by created_at DESC means we need the created_at.
		// To simplify, let's use created_at string as cursor since it's an audit log. Wait, created_at string might lose precision.
		// Let's assume cursor is a valid timestamp string like RFC3339Nano
		query = `SELECT id, user_id, project_id, secret_id, action, ip_address, user_agent, previous_hash, current_hash, created_at FROM secret_access_logs WHERE project_id = $1 AND created_at < $2 ORDER BY created_at DESC LIMIT $3`
		args = []interface{}{projectID, cursor, limit + 1}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var logs []*domain.AuditLog
	for rows.Next() {
		l := &domain.AuditLog{}
		if err := rows.Scan(&l.ID, &l.UserID, &l.ProjectID, &l.SecretID, &l.Action, &l.IPAddress, &l.UserAgent, &l.PreviousHash, &l.CurrentHash, &l.CreatedAt); err != nil {
			return nil, "", err
		}
		logs = append(logs, l)
	}
	if err := rows.Err(); err != nil {
		return nil, "", err
	}

	nextCursor := ""
	if len(logs) > limit {
		nextCursor = logs[limit].CreatedAt.Format(time.RFC3339Nano)
		logs = logs[:limit]
	}

	return logs, nextCursor, nil
}
