package postgres

import (
	"context"
	"database/sql"

	"sitecrawler/newgo/models"
)

type ViewRepo struct {
	db *sql.DB
}

func NewViewRepo(db *sql.DB) *ViewRepo {
	return &ViewRepo{db: db}
}

func (r *ViewRepo) Create(ctx context.Context, v *models.View) error {
	q := `INSERT INTO views (search_keyword_url_id, name, filter_config, created_at, updated_at)
          VALUES ($1,$2,$3,NOW(),NOW()) RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, q, v.SearchKeywordURLID, v.Name, v.FilterConfig).Scan(&v.ID, &v.CreatedAt, &v.UpdatedAt)
}

func (r *ViewRepo) Update(ctx context.Context, v *models.View) error {
	q := `UPDATE views SET name=$2, filter_config=$3, updated_at=NOW() WHERE id=$1`
	_, err := r.db.ExecContext(ctx, q, v.ID, v.Name, v.FilterConfig)
	return err
}

func (r *ViewRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM views WHERE id=$1`, id)
	return err
}

func (r *ViewRepo) Get(ctx context.Context, id int64) (*models.View, error) {
	var v models.View
	q := `SELECT id, search_keyword_url_id, name, filter_config, created_at, updated_at FROM views WHERE id=$1`
	if err := r.db.QueryRowContext(ctx, q, id).Scan(&v.ID, &v.SearchKeywordURLID, &v.Name, &v.FilterConfig, &v.CreatedAt, &v.UpdatedAt); err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *ViewRepo) ListBySKU(ctx context.Context, skuID int64) ([]models.View, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, search_keyword_url_id, name, filter_config, created_at, updated_at FROM views WHERE search_keyword_url_id=$1 ORDER BY created_at ASC`, skuID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.View
	for rows.Next() {
		var v models.View
		if err := rows.Scan(&v.ID, &v.SearchKeywordURLID, &v.Name, &v.FilterConfig, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
