package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"sitecrawler/newgo/models"
)

type AuditRepo struct {
	db *sql.DB
}

func NewAuditRepo(db *sql.DB) *AuditRepo {
	return &AuditRepo{db: db}
}

func (r *AuditRepo) Create(ctx context.Context, ac *models.AuditCheck) error {
	q := `INSERT INTO audit_checks (search_keyword_url_id, name, category, filter_config, created_at, updated_at)
          VALUES ($1,$2,$3,$4,NOW(),NOW()) RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, q, ac.SearchKeywordURLID, ac.Name, ac.Category, ac.FilterConfig).Scan(&ac.ID, &ac.CreatedAt, &ac.UpdatedAt)
}

func (r *AuditRepo) Update(ctx context.Context, ac *models.AuditCheck) error {
	q := `UPDATE audit_checks SET name=$2, category=$3, filter_config=$4, updated_at=NOW() WHERE id=$1`
	_, err := r.db.ExecContext(ctx, q, ac.ID, ac.Name, ac.Category, ac.FilterConfig)
	return err
}

func (r *AuditRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM audit_checks WHERE id=$1`, id)
	return err
}

func (r *AuditRepo) Get(ctx context.Context, id int64) (*models.AuditCheck, error) {
	var ac models.AuditCheck
	q := `SELECT id, search_keyword_url_id, name, category, filter_config, created_at, updated_at FROM audit_checks WHERE id=$1`
	if err := r.db.QueryRowContext(ctx, q, id).Scan(&ac.ID, &ac.SearchKeywordURLID, &ac.Name, &ac.Category, &ac.FilterConfig, &ac.CreatedAt, &ac.UpdatedAt); err != nil {
		return nil, err
	}
	return &ac, nil
}

func (r *AuditRepo) ListBySKU(ctx context.Context, skuID int64) ([]models.AuditCheck, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, search_keyword_url_id, name, category, filter_config, created_at, updated_at FROM audit_checks WHERE search_keyword_url_id=$1 ORDER BY id ASC`, skuID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.AuditCheck
	for rows.Next() {
		var ac models.AuditCheck
		if err := rows.Scan(&ac.ID, &ac.SearchKeywordURLID, &ac.Name, &ac.Category, &ac.FilterConfig, &ac.CreatedAt, &ac.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, ac)
	}
	return out, rows.Err()
}

func (r *AuditRepo) ListByIDsAndSKUs(ctx context.Context, ids, skuIDs []int64) ([]models.AuditCheck, error) {
	if len(ids) == 0 || len(skuIDs) == 0 {
		return nil, nil
	}
	idph := make([]string, len(ids))
	skuPH := make([]string, len(skuIDs))
	args := make([]any, 0, len(ids)+len(skuIDs))
	for i, v := range ids {
		idph[i] = fmt.Sprintf("$%d", i+1)
		args = append(args, v)
	}
	base := len(args)
	for i, v := range skuIDs {
		skuPH[i] = fmt.Sprintf("$%d", base+i+1)
		args = append(args, v)
	}
	q := fmt.Sprintf(`SELECT id, search_keyword_url_id, name, category, filter_config, created_at, updated_at FROM audit_checks WHERE id IN (%s) AND search_keyword_url_id IN (%s) ORDER BY id ASC`, strings.Join(idph, ","), strings.Join(skuPH, ","))
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.AuditCheck
	for rows.Next() {
		var ac models.AuditCheck
		if err := rows.Scan(&ac.ID, &ac.SearchKeywordURLID, &ac.Name, &ac.Category, &ac.FilterConfig, &ac.CreatedAt, &ac.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, ac)
	}
	return out, rows.Err()
}
