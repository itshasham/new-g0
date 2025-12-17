package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
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

func (r *CrawlingSessionRepo) Create(ctx context.Context, session *models.CrawlingSession) error {
	now := time.Now().UTC()
	session.CreatedAt, session.UpdatedAt = now, now
	if session.Status == "" {
		session.Status = "pending"
	}

	q := `INSERT INTO crawling_sessions
		(search_keyword_url_id, url, status, queue, version, options, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id`

	var opt any = session.Options
	if session.Options == nil {
		opt = sql.NullString{}
	}

	return r.db.QueryRowContext(ctx, q,
		session.SearchKeywordURLID, session.URL, session.Status, session.Queue,
		session.Version, opt, session.CreatedAt, session.UpdatedAt,
	).Scan(&session.ID)
}

func (r *CrawlingSessionRepo) PreventInProgress(ctx context.Context, skuID int64) error {
	q := `SELECT 1 FROM crawling_sessions WHERE search_keyword_url_id=$1 AND status IN ('pending','processing') LIMIT 1`
	var one int
	err := r.db.QueryRowContext(ctx, q, skuID).Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return errOrDefault(err, errors.New("session already running"))
}

func (r *CrawlingSessionRepo) GetByID(ctx context.Context, id int64) (*models.CrawlingSession, error) {
	q := `SELECT id, search_keyword_url_id, url, status, queue, version, started_at, ended_at, end_reason, error, created_at, updated_at
		FROM crawling_sessions WHERE id=$1`
	var cs models.CrawlingSession
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&cs.ID, &cs.SearchKeywordURLID, &cs.URL, &cs.Status, &cs.Queue, &cs.Version,
		&cs.StartedAt, &cs.EndedAt, &cs.EndReason, &cs.Error, &cs.CreatedAt, &cs.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("crawling session not found")
		}
		return nil, err
	}
	return &cs, nil
}

func (r *CrawlingSessionRepo) ClaimPending(ctx context.Context, queueID, limit int) ([]models.CrawlingSession, error) {
	if limit <= 0 {
		return nil, nil
	}
	q := `WITH cte AS (
			SELECT id FROM crawling_sessions
			WHERE started_at IS NULL AND queue=$1
			ORDER BY created_at
			FOR UPDATE SKIP LOCKED
			LIMIT $2
		)
		UPDATE crawling_sessions AS cs
		SET started_at = NOW(), status = 'processing', updated_at = NOW()
		FROM cte
		WHERE cs.id = cte.id
		RETURNING cs.id, cs.search_keyword_url_id, cs.url, cs.status, cs.queue, cs.version,
			cs.options, cs.started_at, cs.created_at, cs.updated_at`
	rows, err := r.db.QueryContext(ctx, q, queueID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.CrawlingSession
	for rows.Next() {
		var cs models.CrawlingSession
		var optJSON []byte
		if err := rows.Scan(&cs.ID, &cs.SearchKeywordURLID, &cs.URL, &cs.Status, &cs.Queue,
			&cs.Version, &optJSON, &cs.StartedAt, &cs.CreatedAt, &cs.UpdatedAt); err != nil {
			return nil, err
		}
		if len(optJSON) > 0 {
			_ = json.Unmarshal(optJSON, &cs.Options)
		}
		out = append(out, cs)
	}
	return out, rows.Err()
}

func (r *CrawlingSessionRepo) ClaimStalled(ctx context.Context, queueID int, excludeIDs []int64, limit int) ([]models.CrawlingSession, error) {
	if limit <= 0 {
		return nil, nil
	}
	notIn := ""
	args := []any{queueID}
	if len(excludeIDs) > 0 {
		ph := make([]string, len(excludeIDs))
		for i, id := range excludeIDs {
			ph[i] = fmt.Sprintf("$%d", i+2)
			args = append(args, id)
		}
		notIn = "AND id NOT IN (" + strings.Join(ph, ",") + ")"
	}
	args = append(args, limit)

	q := fmt.Sprintf(`WITH cte AS (
			SELECT id FROM crawling_sessions
			WHERE status='processing' AND queue=$1 %s
			ORDER BY created_at
			FOR UPDATE SKIP LOCKED
			LIMIT $%d
		)
		UPDATE crawling_sessions AS cs
		SET started_at = NOW(), status = 'processing', updated_at = NOW()
		FROM cte
		WHERE cs.id = cte.id
		RETURNING cs.id, cs.search_keyword_url_id, cs.url, cs.status, cs.queue, cs.version,
			cs.started_at, cs.created_at, cs.updated_at`, notIn, len(args))

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.CrawlingSession
	for rows.Next() {
		var cs models.CrawlingSession
		if err := rows.Scan(&cs.ID, &cs.SearchKeywordURLID, &cs.URL, &cs.Status, &cs.Queue,
			&cs.Version, &cs.StartedAt, &cs.CreatedAt, &cs.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, cs)
	}
	return out, rows.Err()
}

func (r *CrawlingSessionRepo) MarkDone(ctx context.Context, id int64, reason string) error {
	q := `UPDATE crawling_sessions SET status='done', end_reason=$2, ended_at=NOW(), updated_at=NOW() WHERE id=$1`
	_, err := r.db.ExecContext(ctx, q, id, reason)
	return err
}

func (r *CrawlingSessionRepo) UpdateSiteInfo(ctx context.Context, id int64, info repository.SiteInfo) error {
	q := `UPDATE crawling_sessions SET ips=$2, dns_servers=$3, aliases=$4, location=$5,
		sitemap=$6, robots=$7, ssl_valid=$8, ssl_valid_until=$9, updated_at=NOW() WHERE id=$1`
	_, err := r.db.ExecContext(ctx, q, id,
		pqTextArray(info.IPs), pqTextArray(info.DNSServers), pqTextArray(info.Aliases),
		info.Location, info.Sitemap, info.Robots, info.SSLValid, info.SSLValidUntil)
	return err
}

func (r *CrawlingSessionRepo) UpdateProgress(ctx context.Context, id int64, d repository.ProgressDelta) error {
	sets := []string{"updated_at=NOW()"}
	if d.IncPages {
		sets = append(sets, "pages_count=COALESCE(pages_count,0)+1")
	}
	if d.InternalURLsDelta != 0 {
		sets = append(sets, fmt.Sprintf("internal_urls_count=COALESCE(internal_urls_count,0)+(%d)", d.InternalURLsDelta))
	}
	if d.IgnoredURLsDelta != 0 {
		sets = append(sets, fmt.Sprintf("ignored_urls_count=COALESCE(ignored_urls_count,0)+(%d)", d.IgnoredURLsDelta))
	}
	if d.ExternalURLsDelta != 0 {
		sets = append(sets, fmt.Sprintf("external_urls_count=COALESCE(external_urls_count,0)+(%d)", d.ExternalURLsDelta))
	}
	q := "UPDATE crawling_sessions SET " + strings.Join(sets, ",") + " WHERE id=$1"
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}

// Helper functions

func pqTextArray(vals []string) any {
	if vals == nil {
		return sql.NullString{}
	}
	b := strings.Builder{}
	b.WriteString("{")
	for i, v := range vals {
		if i > 0 {
			b.WriteString(",")
		}
		esc := strings.ReplaceAll(strings.ReplaceAll(v, `\`, `\\`), `"`, `\"`)
		b.WriteString("\"")
		b.WriteString(esc)
		b.WriteString("\"")
	}
	b.WriteString("}")
	return b.String()
}

func errOrDefault(err error, def error) error {
	if err != nil {
		return err
	}
	return def
}

type CrawlingSessionPageRepo struct {
	db *sql.DB
}

func NewCrawlingSessionPageRepo(db *sql.DB) *CrawlingSessionPageRepo {
	return &CrawlingSessionPageRepo{db: db}
}

func (r *CrawlingSessionPageRepo) List(ctx context.Context, params repository.PageListParams) ([]models.Page, int, error) {
	var args []any
	argIndex := 1

	whereClause := fmt.Sprintf("crawling_session_id = $%d", argIndex)
	args = append(args, params.SessionID)
	argIndex++

	for _, filter := range params.Filters {
		for key, value := range filter {
			whereClause += fmt.Sprintf(" AND %s = $%d", key, argIndex)
			args = append(args, value)
			argIndex++
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

	query := fmt.Sprintf(`SELECT id, crawling_session_id, url, response_code FROM pages WHERE %s ORDER BY %s LIMIT $%d OFFSET $%d`,
		whereClause, orderClause, argIndex, argIndex+1)
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
	var args []any
	argIndex := 1

	whereClause := fmt.Sprintf("p.crawling_session_id = $%d", argIndex)
	args = append(args, params.SessionID)
	argIndex++

	for _, filter := range params.ViewFilters {
		for key, value := range filter {
			whereClause += fmt.Sprintf(" AND p.%s = $%d", key, argIndex)
			args = append(args, value)
			argIndex++
		}
	}

	if params.ComparisonSessionID != nil {
		whereClause += fmt.Sprintf(" OR p.crawling_session_id = $%d", argIndex)
		args = append(args, *params.ComparisonSessionID)
	}

	query := fmt.Sprintf(`SELECT c.id, c.name, p.id, p.crawling_session_id, p.url, p.response_code
		FROM audit_checks c LEFT JOIN pages p ON p.crawling_session_id = $1
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
		WHERE p.id = $1 AND cs.search_keyword_url_id = $2`

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
	rows, err := r.db.QueryContext(ctx, `SELECT id, page_id, url FROM page_images WHERE page_id = $1 LIMIT $2`, pageID, limit)
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
		WHERE pl.source_page_id = $1 AND p.response_code >= 400 LIMIT $2`

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
		WHERE pl.target_page_id = $1 LIMIT $2`

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
	whereClause, args, err := buildPagesWherePostgres(params.SessionID, params.Prefilters, params.Filters)
	if err != nil {
		return nil, err
	}

	q := fmt.Sprintf(`SELECT
		COUNT(*) AS total,
		COUNT(*) FILTER (WHERE (og_title = '' OR og_title IS NULL OR og_description = '' OR og_description IS NULL)) AS warning,
		COUNT(*) FILTER (WHERE (response_code >= 400 AND response_code <= 599)) AS error,
		COUNT(*) FILTER (WHERE (response_code >= 200 AND response_code <= 299 AND redirect_code IS NULL)) AS ok,
		COUNT(*) FILTER (WHERE (redirect_code IN ('301','302','307','308'))) AS redirection,
		COUNT(*) FILTER (WHERE (depth = 1)) AS level1,
		COUNT(*) FILTER (WHERE (depth = 2)) AS level2,
		COUNT(*) FILTER (WHERE (depth = 3)) AS level3,
		COUNT(*) FILTER (WHERE (depth = 4)) AS level4,
		COUNT(*) FILTER (WHERE (response_code >= 200 AND response_code < 300)) AS success_pages,
		COUNT(*) FILTER (WHERE (response_code >= 300 AND response_code < 400)) AS redirect_pages,
		COUNT(*) FILTER (WHERE (response_code >= 400 AND response_code < 500)) AS client_error_pages,
		COUNT(*) FILTER (WHERE (response_code >= 500)) AS server_error_pages
		FROM pages WHERE %s`, whereClause)

	var rTotal, rWarning, rError, rOK, rRedirection int
	var rL1, rL2, rL3, rL4 int
	var rSuccess, rResp3xx, rClientErr, rServerErr int
	if err := r.db.QueryRowContext(ctx, q, args...).Scan(
		&rTotal, &rWarning, &rError, &rOK, &rRedirection,
		&rL1, &rL2, &rL3, &rL4,
		&rSuccess, &rResp3xx, &rClientErr, &rServerErr,
	); err != nil {
		return nil, err
	}

	result := map[string]any{
		"total": rTotal, "warning": rWarning, "error": rError, "ok": rOK, "redirection": rRedirection,
		"level1": rL1, "level2": rL2, "level3": rL3, "level4": rL4,
		"total_pages": rTotal, "success_pages": rSuccess, "redirect_pages": rResp3xx,
		"client_error_pages": rClientErr, "server_error_pages": rServerErr,
	}

	problematicCount, err := r.fetchProblematicCount(ctx, whereClause, args, params.SessionID)
	if err == nil {
		result["problematic"] = problematicCount
		if rTotal > 0 && problematicCount > 0 {
			result["site_health"] = 100 - (problematicCount * 100 / rTotal)
		} else if rTotal > 0 {
			result["site_health"] = 100
		} else {
			result["site_health"] = 0
		}
	}

	if params.ComparisonSessionID != nil {
		cWhere, cArgs, err := buildPagesWherePostgres(*params.ComparisonSessionID, params.Prefilters, params.Filters)
		if err == nil {
			cq := fmt.Sprintf(`SELECT
				COUNT(*) AS total,
				COUNT(*) FILTER (WHERE (og_title = '' OR og_title IS NULL OR og_description = '' OR og_description IS NULL)) AS warning,
				COUNT(*) FILTER (WHERE (response_code >= 400 AND response_code <= 599)) AS error,
				COUNT(*) FILTER (WHERE (response_code >= 200 AND response_code <= 299 AND redirect_code IS NULL)) AS ok,
				COUNT(*) FILTER (WHERE (redirect_code IN ('301','302','307','308'))) AS redirection
				FROM pages WHERE %s`, cWhere)
			var cTotal, cWarning, cError, cOK, cRedirection int
			if err := r.db.QueryRowContext(ctx, cq, cArgs...).Scan(&cTotal, &cWarning, &cError, &cOK, &cRedirection); err == nil {
				result["comparison"] = map[string]any{
					"total": cTotal, "warning": cWarning, "error": cError, "ok": cOK, "redirection": cRedirection,
				}
				result["changes"] = map[string]int{
					"total": rTotal - cTotal, "warning": rWarning - cWarning, "error": rError - cError, "ok": rOK - cOK, "redirection": rRedirection - cRedirection,
				}
			}
		}
	}

	return result, nil
}

func (r *StatsRepo) fetchProblematicCount(ctx context.Context, baseWhere string, baseArgs []any, sessionID int64) (int, error) {
	var skuID int64
	if err := r.db.QueryRowContext(ctx, `SELECT search_keyword_url_id FROM crawling_sessions WHERE id = $1`, sessionID).Scan(&skuID); err != nil {
		return 0, err
	}

	rows, err := r.db.QueryContext(ctx, `SELECT filter_config FROM audit_checks WHERE search_keyword_url_id = $1 AND category = 'problematic'`, skuID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var rawConfigs [][]byte
	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			return 0, err
		}
		if len(raw) > 0 {
			rawConfigs = append(rawConfigs, raw)
		}
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}

	probClause, probArgs, err := buildProblematicClausePostgres(rawConfigs)
	if err != nil || probClause == "" {
		return 0, err
	}

	renumbered := renumberPostgresPlaceholders(probClause, len(baseArgs)+1)
	args := append(append([]any{}, baseArgs...), probArgs...)
	q := fmt.Sprintf("SELECT COUNT(*) FROM pages WHERE %s AND (%s)", baseWhere, renumbered)
	var count int
	if err := r.db.QueryRowContext(ctx, q, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func renumberPostgresPlaceholders(clause string, start int) string {
	// buildProblematicClausePostgres uses $1..$N, but we need to append after existing args.
	re := regexp.MustCompile(`\$(\d+)`)
	return re.ReplaceAllStringFunc(clause, func(m string) string {
		n, err := strconv.Atoi(strings.TrimPrefix(m, "$"))
		if err != nil {
			return m
		}
		return fmt.Sprintf("$%d", start+n-1)
	})
}
