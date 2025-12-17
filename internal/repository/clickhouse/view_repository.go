package clickhouse

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"sitecrawler/newgo/models"
)

type ViewRepo struct {
	db *sql.DB
}

func NewViewRepo(db *sql.DB) *ViewRepo {
	return &ViewRepo{db: db}
}

func (r *ViewRepo) Create(ctx context.Context, v *models.View) error {
	now := time.Now().UTC()
	v.CreatedAt, v.UpdatedAt = now, now

	filterJSON, err := json.Marshal(v.FilterConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal filter config: %w", err)
	}

	q := `INSERT INTO views (search_keyword_url_id, name, filter_config, created_at, updated_at)
	      VALUES (?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, q, v.SearchKeywordURLID, v.Name, string(filterJSON), v.CreatedAt, v.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		err = r.db.QueryRowContext(ctx, "SELECT max(id) + 1 FROM views").Scan(&id)
		if err != nil {
			return fmt.Errorf("failed to get view ID: %w", err)
		}
	}
	v.ID = id
	return nil
}

func (r *ViewRepo) Update(ctx context.Context, v *models.View) error {
	filterJSON, err := json.Marshal(v.FilterConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal filter config: %w", err)
	}

	q := `ALTER TABLE views UPDATE name = ?, filter_config = ?, updated_at = ? WHERE id = ?`
	_, err = r.db.ExecContext(ctx, q, v.Name, string(filterJSON), time.Now().UTC(), v.ID)
	return err
}

func (r *ViewRepo) Delete(ctx context.Context, id int64) error {
	q := `ALTER TABLE views DELETE WHERE id = ?`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}

func (r *ViewRepo) Get(ctx context.Context, id int64) (*models.View, error) {
	q := `SELECT id, search_keyword_url_id, name, filter_config, created_at, updated_at 
	      FROM views WHERE id = ?`

	var v models.View
	var filterJSON string

	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&v.ID, &v.SearchKeywordURLID, &v.Name, &filterJSON, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if filterJSON != "" {
		if err := json.Unmarshal([]byte(filterJSON), &v.FilterConfig); err != nil {
			return nil, fmt.Errorf("failed to unmarshal filter config: %w", err)
		}
	}

	return &v, nil
}

func (r *ViewRepo) ListBySKU(ctx context.Context, skuID int64) ([]models.View, error) {
	q := `SELECT id, search_keyword_url_id, name, filter_config, created_at, updated_at 
	      FROM views WHERE search_keyword_url_id = ? ORDER BY created_at ASC`

	rows, err := r.db.QueryContext(ctx, q, skuID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []models.View
	for rows.Next() {
		var v models.View
		var filterJSON string

		err := rows.Scan(&v.ID, &v.SearchKeywordURLID, &v.Name, &filterJSON, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if filterJSON != "" {
			if err := json.Unmarshal([]byte(filterJSON), &v.FilterConfig); err != nil {
				return nil, fmt.Errorf("failed to unmarshal filter config: %w", err)
			}
		}

		views = append(views, v)
	}

	return views, rows.Err()
}
