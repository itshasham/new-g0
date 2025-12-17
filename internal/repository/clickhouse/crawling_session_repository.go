package clickhouse

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"sitecrawler/newgo/internal/repository"
	"sitecrawler/newgo/models"
)

type CrawlingSessionRepo struct {
	db *sql.DB
}

func NewCrawlingSessionRepo(db *sql.DB) *CrawlingSessionRepo {
	return &CrawlingSessionRepo{db: db}
}

func (r *CrawlingSessionRepo) Create(ctx context.Context, cs *models.CrawlingSession) error {
	now := time.Now().UTC()
	cs.CreatedAt, cs.UpdatedAt = now, now
	if cs.Status == "" {
		cs.Status = "pending"
	}

	// Serialize options to JSON
	var optJSON []byte
	var err error
	if cs.Options != nil {
		optJSON, err = json.Marshal(cs.Options)
		if err != nil {
			return fmt.Errorf("failed to marshal options: %w", err)
		}
	}

	// ClickHouse INSERT
	q := `INSERT INTO crawling_sessions (
		search_keyword_url_id, url, status, queue, version, options, 
		created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, q,
		cs.SearchKeywordURLID, cs.URL, cs.Status, cs.Queue, cs.Version,
		string(optJSON), cs.CreatedAt, cs.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// ClickHouse doesn't auto-generate IDs like Postgres SERIAL
	// We need to either use a separate sequence or generate IDs client-side
	id, err := result.LastInsertId()
	if err != nil {
		// Fallback: query max ID + 1 (not ideal for concurrency)
		err = r.db.QueryRowContext(ctx, "SELECT max(id) + 1 FROM crawling_sessions").Scan(&id)
		if err != nil {
			return fmt.Errorf("failed to get session ID: %w", err)
		}
	}
	cs.ID = id
	return nil
}

func (r *CrawlingSessionRepo) PreventInProgress(ctx context.Context, skuID int64) error {
	// ClickHouse doesn't support SELECT FOR UPDATE, so we use a simple check
	var count int
	q := `SELECT count(*) FROM crawling_sessions 
	      WHERE search_keyword_url_id = ? AND status IN ('pending', 'processing')`
	err := r.db.QueryRowContext(ctx, q, skuID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("session already running")
	}
	return nil
}

func (r *CrawlingSessionRepo) GetByID(ctx context.Context, id int64) (*models.CrawlingSession, error) {
	q := `SELECT id, search_keyword_url_id, url, status, queue, version,
	             started_at, ended_at, end_reason, error, created_at, updated_at
	      FROM crawling_sessions WHERE id = ?`

	var cs models.CrawlingSession
	var startedAt, endedAt int64
	var endReason, errStr sql.NullString

	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&cs.ID, &cs.SearchKeywordURLID, &cs.URL, &cs.Status, &cs.Queue, &cs.Version,
		&startedAt, &endedAt, &endReason, &errStr, &cs.CreatedAt, &cs.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("crawling session not found")
		}
		return nil, err
	}

	if startedAt > 0 {
		t := time.Unix(startedAt, 0)
		cs.StartedAt = &t
	}
	if endedAt > 0 {
		t := time.Unix(endedAt, 0)
		cs.EndedAt = &t
	}
	if endReason.Valid {
		cs.EndReason = endReason.String
	}
	if errStr.Valid {
		cs.Error = errStr.String
	}

	return &cs, nil
}

func (r *CrawlingSessionRepo) ClaimPending(ctx context.Context, queueID, limit int) ([]models.CrawlingSession, error) {
	// ClickHouse doesn't support UPDATE...RETURNING or row-level locking
	// This is a simplified implementation - in production, you'd need a different strategy
	// such as using a separate coordination service (Redis, etcd) or optimistic locking

	if limit <= 0 {
		return nil, nil
	}

	q := `SELECT id, search_keyword_url_id, url, status, queue, version, options, 
	             started_at, created_at, updated_at
	      FROM crawling_sessions
	      WHERE started_at = 0 AND queue = ?
	      ORDER BY created_at
	      LIMIT ?`

	rows, err := r.db.QueryContext(ctx, q, queueID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.CrawlingSession
	var ids []int64

	for rows.Next() {
		var cs models.CrawlingSession
		var optJSON string
		var startedAt int64

		err := rows.Scan(&cs.ID, &cs.SearchKeywordURLID, &cs.URL, &cs.Status,
			&cs.Queue, &cs.Version, &optJSON, &startedAt, &cs.CreatedAt, &cs.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if optJSON != "" {
			if err := json.Unmarshal([]byte(optJSON), &cs.Options); err != nil {
				return nil, err
			}
		}

		sessions = append(sessions, cs)
		ids = append(ids, cs.ID)
	}

	if len(ids) == 0 {
		return nil, nil
	}

	// Update claimed sessions (in ClickHouse, uses ALTER TABLE UPDATE)
	// Note: This is not atomic and may have race conditions in production
	now := time.Now().UTC()
	for _, id := range ids {
		_, err := r.db.ExecContext(ctx,
			`ALTER TABLE crawling_sessions UPDATE started_at = ?, status = 'processing', updated_at = ? WHERE id = ?`,
			now.Unix(), now, id)
		if err != nil {
			return nil, err
		}
	}

	return sessions, nil
}

func (r *CrawlingSessionRepo) ClaimStalled(ctx context.Context, queueID int, excludeIDs []int64, limit int) ([]models.CrawlingSession, error) {
	// ClickHouse doesn't support FOR UPDATE / SKIP LOCKED
	// This is a simplified implementation
	return nil, fmt.Errorf("ClaimStalled not yet implemented for ClickHouse")
}

func (r *CrawlingSessionRepo) MarkDone(ctx context.Context, id int64, reason string) error {
	now := time.Now().UTC()
	q := `ALTER TABLE crawling_sessions UPDATE
	      status = 'done', ended_at = ?, end_reason = ?, updated_at = ?
	      WHERE id = ?`
	_, err := r.db.ExecContext(ctx, q, now.Unix(), reason, now, id)
	return err
}

func (r *CrawlingSessionRepo) UpdateSiteInfo(ctx context.Context, id int64, info repository.SiteInfo) error {
	// Convert slices to JSON arrays for ClickHouse
	ipsJSON, _ := json.Marshal(info.IPs)
	dnsJSON, _ := json.Marshal(info.DNSServers)
	aliasesJSON, _ := json.Marshal(info.Aliases)

	var sslValidUntil int64
	if info.SSLValidUntil != nil {
		sslValidUntil = info.SSLValidUntil.Unix()
	}

	q := `ALTER TABLE crawling_sessions UPDATE
	      ips = ?, dns_servers = ?, aliases = ?, location = ?,
	      sitemap = ?, robots = ?, ssl_valid = ?, ssl_valid_until = ?, updated_at = ?
	      WHERE id = ?`

	_, err := r.db.ExecContext(ctx, q,
		string(ipsJSON), string(dnsJSON), string(aliasesJSON), info.Location,
		info.Sitemap, info.Robots, info.SSLValid, sslValidUntil, time.Now().UTC(), id)
	return err
}

func (r *CrawlingSessionRepo) UpdateProgress(ctx context.Context, id int64, d repository.ProgressDelta) error {
	// ClickHouse doesn't support traditional UPDATE with increments
	// We need to use ALTER TABLE UPDATE with explicit values
	// This requires reading current values first (not ideal)
	return fmt.Errorf("UpdateProgress not yet fully implemented for ClickHouse")
}

type CrawlingSessionPageRepo struct {
	db *sql.DB
}

func NewCrawlingSessionPageRepo(db *sql.DB) *CrawlingSessionPageRepo {
	return &CrawlingSessionPageRepo{db: db}
}

func (r *CrawlingSessionPageRepo) List(ctx context.Context, params repository.PageListParams) ([]models.Page, int, error) {
	whereClause := "crawling_session_id = ?"
	args := []any{params.SessionID}

	for _, filter := range params.Filters {
		for key, value := range filter {
			whereClause += fmt.Sprintf(" AND %s = ?", key)
			args = append(args, value)
		}
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM pages WHERE %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	orderClause := "id ASC"
	if params.Sort != "" {
		direction := "ASC"
		if strings.ToUpper(params.Direction) == "DESC" {
			direction = "DESC"
		}
		orderClause = fmt.Sprintf("%s %s", params.Sort, direction)
	}

	limit := params.PageLimit
	if limit <= 0 {
		limit = 20
	}
	offset := 0
	if params.Page > 1 {
		offset = (params.Page - 1) * limit
	}

	query := fmt.Sprintf(`SELECT id, crawling_session_id, url, response_code FROM pages WHERE %s ORDER BY %s LIMIT ? OFFSET ?`, whereClause, orderClause)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var pages []models.Page
	for rows.Next() {
		var p models.Page
		if err := rows.Scan(&p.ID, &p.CrawlingSessionID, &p.URL, &p.ResponseCode); err != nil {
			return nil, 0, err
		}
		pages = append(pages, p)
	}
	return pages, total, rows.Err()
}

type CrawlingSessionCheckRepo struct {
	db *sql.DB
}

func NewCrawlingSessionCheckRepo(db *sql.DB) *CrawlingSessionCheckRepo {
	return &CrawlingSessionCheckRepo{db: db}
}

func (r *CrawlingSessionCheckRepo) ChecksWithPages(ctx context.Context, params repository.ChecksWithPagesParams) ([]models.CheckWithPages, error) {
	whereClause := "p.crawling_session_id = ?"
	args := []any{params.SessionID}

	for _, filter := range params.ViewFilters {
		for key, value := range filter {
			whereClause += fmt.Sprintf(" AND p.%s = ?", key)
			args = append(args, value)
		}
	}

	if params.ComparisonSessionID != nil {
		whereClause += " OR p.crawling_session_id = ?"
		args = append(args, *params.ComparisonSessionID)
	}

	query := fmt.Sprintf(`SELECT c.id, c.name, p.id, p.crawling_session_id, p.url, p.response_code
		FROM audit_checks c LEFT JOIN pages p ON p.crawling_session_id = ?
		WHERE %s ORDER BY c.id ASC, p.id ASC`, whereClause)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checksMap := make(map[int64]*models.CheckWithPages)
	var orderedIDs []int64

	for rows.Next() {
		var checkID int64
		var checkName string
		var pageID, pageCrawlingSessionID sql.NullInt64
		var pageURL sql.NullString
		var pageResponseCode sql.NullInt64

		if err := rows.Scan(&checkID, &checkName, &pageID, &pageCrawlingSessionID, &pageURL, &pageResponseCode); err != nil {
			return nil, err
		}

		check, exists := checksMap[checkID]
		if !exists {
			check = &models.CheckWithPages{ID: checkID, Name: checkName, Pages: []models.Page{}}
			checksMap[checkID] = check
			orderedIDs = append(orderedIDs, checkID)
		}

		if pageID.Valid && len(check.Pages) < params.PageLimitPerCheck {
			check.Pages = append(check.Pages, models.Page{
				ID: pageID.Int64, CrawlingSessionID: pageCrawlingSessionID.Int64,
				URL: pageURL.String, ResponseCode: int(pageResponseCode.Int64),
			})
		}
	}

	var result []models.CheckWithPages
	for _, id := range orderedIDs {
		result = append(result, *checksMap[id])
	}
	return result, rows.Err()
}

type PageDetailsRepo struct {
	db *sql.DB
}

func NewPageDetailsRepo(db *sql.DB) *PageDetailsRepo {
	return &PageDetailsRepo{db: db}
}

func (r *PageDetailsRepo) GetPageByIDAndSKU(ctx context.Context, pageID, skuID int64) (*models.Page, error) {
	q := `SELECT p.id, p.crawling_session_id, p.url, p.response_code
		FROM pages p JOIN crawling_sessions cs ON p.crawling_session_id = cs.id
		WHERE p.id = ? AND cs.search_keyword_url_id = ?`

	var p models.Page
	err := r.db.QueryRowContext(ctx, q, pageID, skuID).Scan(&p.ID, &p.CrawlingSessionID, &p.URL, &p.ResponseCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("page not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *PageDetailsRepo) GetPageImages(ctx context.Context, pageID int64, limit int) ([]models.PageImage, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, page_id, url FROM page_images WHERE page_id = ? LIMIT ?`, pageID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.PageImage
	for rows.Next() {
		var img models.PageImage
		if err := rows.Scan(&img.ID, &img.PageID, &img.URL); err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, rows.Err()
}

func (r *PageDetailsRepo) GetBrokenTargetsFrom(ctx context.Context, pageID int64, limit int) ([]models.Page, error) {
	q := `SELECT p.id, p.crawling_session_id, p.url, p.response_code
		FROM pages p JOIN page_links pl ON pl.target_page_id = p.id
		WHERE pl.source_page_id = ? AND p.response_code >= 400 LIMIT ?`

	rows, err := r.db.QueryContext(ctx, q, pageID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []models.Page
	for rows.Next() {
		var p models.Page
		if err := rows.Scan(&p.ID, &p.CrawlingSessionID, &p.URL, &p.ResponseCode); err != nil {
			return nil, err
		}
		pages = append(pages, p)
	}
	return pages, rows.Err()
}

func (r *PageDetailsRepo) GetReferrersToBroken(ctx context.Context, pageID int64, limit int) ([]models.Page, error) {
	q := `SELECT p.id, p.crawling_session_id, p.url, p.response_code
		FROM pages p JOIN page_links pl ON pl.source_page_id = p.id
		WHERE pl.target_page_id = ? LIMIT ?`

	rows, err := r.db.QueryContext(ctx, q, pageID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []models.Page
	for rows.Next() {
		var p models.Page
		if err := rows.Scan(&p.ID, &p.CrawlingSessionID, &p.URL, &p.ResponseCode); err != nil {
			return nil, err
		}
		pages = append(pages, p)
	}
	return pages, rows.Err()
}

type StatsRepo struct {
	db *sql.DB
}

func NewStatsRepo(db *sql.DB) *StatsRepo {
	return &StatsRepo{db: db}
}

func (r *StatsRepo) Fetch(ctx context.Context, params repository.StatsQueryParams) (map[string]any, error) {
	whereClause := "crawling_session_id = ?"
	args := []any{params.SessionID}

	for _, filter := range params.Prefilters {
		for key, value := range filter {
			whereClause += fmt.Sprintf(" AND %s = ?", key)
			args = append(args, value)
		}
	}

	for _, filter := range params.Filters {
		for key, value := range filter {
			whereClause += fmt.Sprintf(" AND %s = ?", key)
			args = append(args, value)
		}
	}

	query := fmt.Sprintf(`SELECT COUNT(*) as total_pages,
		countIf(response_code >= 200 AND response_code < 300) as success_pages,
		countIf(response_code >= 300 AND response_code < 400) as redirect_pages,
		countIf(response_code >= 400 AND response_code < 500) as client_error_pages,
		countIf(response_code >= 500) as server_error_pages
		FROM pages WHERE %s`, whereClause)

	var totalPages, successPages, redirectPages, clientErrorPages, serverErrorPages int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&totalPages, &successPages, &redirectPages, &clientErrorPages, &serverErrorPages)
	if err != nil {
		return nil, err
	}

	result := map[string]any{
		"total_pages": totalPages, "success_pages": successPages, "redirect_pages": redirectPages,
		"client_error_pages": clientErrorPages, "server_error_pages": serverErrorPages,
	}

	if params.ComparisonSessionID != nil {
		compQuery := `SELECT COUNT(*) as total_pages,
			countIf(response_code >= 200 AND response_code < 300) as success_pages,
			countIf(response_code >= 300 AND response_code < 400) as redirect_pages,
			countIf(response_code >= 400 AND response_code < 500) as client_error_pages,
			countIf(response_code >= 500) as server_error_pages
			FROM pages WHERE crawling_session_id = ?`

		var compTotal, compSuccess, compRedirect, compClientError, compServerError int
		err := r.db.QueryRowContext(ctx, compQuery, *params.ComparisonSessionID).Scan(&compTotal, &compSuccess, &compRedirect, &compClientError, &compServerError)
		if err == nil {
			result["comparison"] = map[string]any{
				"total_pages": compTotal, "success_pages": compSuccess, "redirect_pages": compRedirect,
				"client_error_pages": compClientError, "server_error_pages": compServerError,
			}
		}
	}

	return result, nil
}
