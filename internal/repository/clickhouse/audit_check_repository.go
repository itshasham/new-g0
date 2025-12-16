package clickhouse

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"sitecrawler/newgo/models"
)

type AuditRepo struct {
	db *sql.DB
}

func NewAuditRepo(db *sql.DB) *AuditRepo {
	return &AuditRepo{db: db}
}

func (r *AuditRepo) Create(ctx context.Context, ac *models.AuditCheck) error {
	now := time.Now().UTC()
	ac.CreatedAt, ac.UpdatedAt = now, now

	filterJSON, err := json.Marshal(ac.FilterConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal filter config: %w", err)
	}

	q := `INSERT INTO audit_checks (search_keyword_url_id, name, category, filter_config, created_at, updated_at)
	      VALUES (?, ?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, q, ac.SearchKeywordURLID, ac.Name, ac.Category, string(filterJSON), ac.CreatedAt, ac.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		err = r.db.QueryRowContext(ctx, "SELECT max(id) + 1 FROM audit_checks").Scan(&id)
		if err != nil {
			return fmt.Errorf("failed to get audit check ID: %w", err)
		}
	}
	ac.ID = id
	return nil
}

func (r *AuditRepo) Update(ctx context.Context, ac *models.AuditCheck) error {
	filterJSON, err := json.Marshal(ac.FilterConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal filter config: %w", err)
	}

	q := `ALTER TABLE audit_checks UPDATE name = ?, category = ?, filter_config = ?, updated_at = ? WHERE id = ?`
	_, err = r.db.ExecContext(ctx, q, ac.Name, ac.Category, string(filterJSON), time.Now().UTC(), ac.ID)
	return err
}

func (r *AuditRepo) Delete(ctx context.Context, id int64) error {
	q := `ALTER TABLE audit_checks DELETE WHERE id = ?`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}

func (r *AuditRepo) Get(ctx context.Context, id int64) (*models.AuditCheck, error) {
	q := `SELECT id, search_keyword_url_id, name, category, filter_config, created_at, updated_at 
	      FROM audit_checks WHERE id = ?`

	var ac models.AuditCheck
	var filterJSON string

	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&ac.ID, &ac.SearchKeywordURLID, &ac.Name, &ac.Category, &filterJSON, &ac.CreatedAt, &ac.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if filterJSON != "" {
		if err := json.Unmarshal([]byte(filterJSON), &ac.FilterConfig); err != nil {
			return nil, fmt.Errorf("failed to unmarshal filter config: %w", err)
		}
	}

	return &ac, nil
}

func (r *AuditRepo) ListBySKU(ctx context.Context, skuID int64) ([]models.AuditCheck, error) {
	q := `SELECT id, search_keyword_url_id, name, category, filter_config, created_at, updated_at 
	      FROM audit_checks WHERE search_keyword_url_id = ? ORDER BY created_at ASC`

	rows, err := r.db.QueryContext(ctx, q, skuID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checks []models.AuditCheck
	for rows.Next() {
		var ac models.AuditCheck
		var filterJSON string

		err := rows.Scan(&ac.ID, &ac.SearchKeywordURLID, &ac.Name, &ac.Category, &filterJSON, &ac.CreatedAt, &ac.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if filterJSON != "" {
			if err := json.Unmarshal([]byte(filterJSON), &ac.FilterConfig); err != nil {
				return nil, fmt.Errorf("failed to unmarshal filter config: %w", err)
			}
		}

		checks = append(checks, ac)
	}

	return checks, rows.Err()
}

func (r *AuditRepo) ListByIDsAndSKUs(ctx context.Context, ids, skuIDs []int64) ([]models.AuditCheck, error) {
	_ = ctx
	_ = ids
	_ = skuIDs
	return nil, fmt.Errorf("ListByIDsAndSKUs not yet implemented for ClickHouse")
}
