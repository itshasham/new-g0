package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"sitecrawler/newgo/models"
)

// ErrCrawlingSessionNotFound is returned when a session cannot be located.
var ErrCrawlingSessionNotFound = errors.New("crawling session not found")
var ErrPageNotFound = errors.New("page not found")

// CrawlingSessionRepository describes the data access needed by the service layer.
type CrawlingSessionRepository interface {
	PreventInProgress(ctx context.Context, skuID int64) error
	Create(ctx context.Context, session *models.CrawlingSession) error
	GetByID(ctx context.Context, id int64) (*models.CrawlingSession, error)
	ClaimPending(ctx context.Context, queueID, limit int) ([]models.CrawlingSession, error)
	ClaimStalled(ctx context.Context, queueID int, excludeIDs []int64, limit int) ([]models.CrawlingSession, error)
	MarkDone(ctx context.Context, id int64, reason string) error
	UpdateSiteInfo(ctx context.Context, id int64, info SiteInfo) error
	UpdateProgress(ctx context.Context, id int64, d ProgressDelta) error
}

// SiteInfo contains site metadata from crawling
type SiteInfo struct {
	IPs           []string
	DNSServers    []string
	Aliases       []string
	Location      string
	Sitemap       bool
	Robots        bool
	SSLValid      bool
	SSLValidUntil *time.Time
}

// ProgressDelta contains incremental progress updates
type ProgressDelta struct {
	IncPages          bool
	InternalURLsDelta int
	IgnoredURLsDelta  int
	ExternalURLsDelta int
}

type PageListParams struct {
	SessionID int64
	Filters   []map[string]any
	Sort      string
	Direction string
	Page      int
	PageLimit int
}

type ChecksWithPagesParams struct {
	SessionID           int64
	ComparisonSessionID *int64
	ViewFilters         []map[string]any
	PageLimitPerCheck   int
}

type CrawlingSessionPageRepository interface {
	List(ctx context.Context, params PageListParams) ([]models.Page, int, error)
}

type CrawlingSessionCheckRepository interface {
	ChecksWithPages(ctx context.Context, params ChecksWithPagesParams) ([]models.CheckWithPages, error)
}

type PageDetailsRepository interface {
	GetPageByIDAndSKU(ctx context.Context, pageID, skuID int64) (*models.Page, error)
	GetPageImages(ctx context.Context, pageID int64, limit int) ([]models.PageImage, error)
	GetBrokenTargetsFrom(ctx context.Context, pageID int64, limit int) ([]models.Page, error)
	GetReferrersToBroken(ctx context.Context, pageID int64, limit int) ([]models.Page, error)
}

type StatsQueryParams struct {
	SessionID           int64
	ComparisonSessionID *int64
	Filters             []map[string]any
	Prefilters          []map[string]any
}

type StatsRepository interface {
	Fetch(ctx context.Context, params StatsQueryParams) (map[string]any, error)
}

// InMemoryCrawlingSessionRepository is a temporary implementation that stores sessions
// in memory so the new architecture can be exercised without database dependencies.
type InMemoryCrawlingSessionRepository struct {
	mu        sync.Mutex
	nextID    int64
	activeSKU map[int64]struct{}
	items     map[int64]*models.CrawlingSession
}

func NewInMemoryCrawlingSessionRepository() *InMemoryCrawlingSessionRepository {
	return &InMemoryCrawlingSessionRepository{
		nextID:    1,
		activeSKU: make(map[int64]struct{}),
		items:     make(map[int64]*models.CrawlingSession),
	}
}

func (r *InMemoryCrawlingSessionRepository) PreventInProgress(ctx context.Context, skuID int64) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.activeSKU[skuID]; exists {
		return errors.New("crawling session already in progress")
	}
	r.activeSKU[skuID] = struct{}{}
	return nil
}

func (r *InMemoryCrawlingSessionRepository) Create(ctx context.Context, session *models.CrawlingSession) error {
	if session == nil {
		return errors.New("session must not be nil")
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	now := time.Now().UTC()

	r.mu.Lock()
	defer r.mu.Unlock()
	session.ID = r.nextID
	r.nextID++
	session.CreatedAt = now
	session.UpdatedAt = now
	// mark SKU as no longer active once persisted
	delete(r.activeSKU, session.SearchKeywordURLID)
	r.items[session.ID] = session
	return nil
}

func (r *InMemoryCrawlingSessionRepository) GetByID(ctx context.Context, id int64) (*models.CrawlingSession, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	session, ok := r.items[id]
	if !ok {
		return nil, ErrCrawlingSessionNotFound
	}

	copied := *session
	return &copied, nil
}

func (r *InMemoryCrawlingSessionRepository) ClaimPending(ctx context.Context, queueID, limit int) ([]models.CrawlingSession, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []models.CrawlingSession
	for _, s := range r.items {
		if s.Queue == queueID && s.Status == "pending" && len(out) < limit {
			s.Status = "processing"
			now := time.Now().UTC()
			s.StartedAt = &now
			out = append(out, *s)
		}
	}
	return out, nil
}

func (r *InMemoryCrawlingSessionRepository) ClaimStalled(ctx context.Context, queueID int, excludeIDs []int64, limit int) ([]models.CrawlingSession, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	exclude := make(map[int64]bool)
	for _, id := range excludeIDs {
		exclude[id] = true
	}
	var out []models.CrawlingSession
	for _, s := range r.items {
		if s.Queue == queueID && s.Status == "processing" && !exclude[s.ID] && len(out) < limit {
			now := time.Now().UTC()
			s.StartedAt = &now
			out = append(out, *s)
		}
	}
	return out, nil
}

func (r *InMemoryCrawlingSessionRepository) MarkDone(ctx context.Context, id int64, reason string) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	if s, ok := r.items[id]; ok {
		s.Status = "done"
		s.EndReason = reason
		now := time.Now().UTC()
		s.EndedAt = &now
	}
	return nil
}

func (r *InMemoryCrawlingSessionRepository) UpdateSiteInfo(ctx context.Context, id int64, info SiteInfo) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	if s, ok := r.items[id]; ok {
		s.IPs = info.IPs
		s.DNSServers = info.DNSServers
		s.Aliases = info.Aliases
		s.Location = info.Location
		s.Sitemap = info.Sitemap
		s.Robots = info.Robots
		s.SSLValid = info.SSLValid
		s.SSLValidUntil = info.SSLValidUntil
	}
	return nil
}

func (r *InMemoryCrawlingSessionRepository) UpdateProgress(ctx context.Context, id int64, d ProgressDelta) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	if s, ok := r.items[id]; ok {
		if d.IncPages {
			s.PagesCount++
		}
		s.InternalURLsCount += d.InternalURLsDelta
		s.IgnoredURLsCount += d.IgnoredURLsDelta
		s.ExternalURLsCount += d.ExternalURLsDelta
	}
	return nil
}

type NoopCrawlingSessionPageRepository struct{}

func NewNoopCrawlingSessionPageRepository() *NoopCrawlingSessionPageRepository {
	return &NoopCrawlingSessionPageRepository{}
}

func (r *NoopCrawlingSessionPageRepository) List(ctx context.Context, params PageListParams) ([]models.Page, int, error) {
	_ = ctx
	_ = params
	return nil, 0, nil
}

type NoopCrawlingSessionCheckRepository struct{}

func NewNoopCrawlingSessionCheckRepository() *NoopCrawlingSessionCheckRepository {
	return &NoopCrawlingSessionCheckRepository{}
}

func (r *NoopCrawlingSessionCheckRepository) ChecksWithPages(ctx context.Context, params ChecksWithPagesParams) ([]models.CheckWithPages, error) {
	_ = ctx
	_ = params
	return nil, nil
}

type NoopPageDetailsRepository struct{}

func NewNoopPageDetailsRepository() *NoopPageDetailsRepository {
	return &NoopPageDetailsRepository{}
}

func (r *NoopPageDetailsRepository) GetPageByIDAndSKU(ctx context.Context, pageID, skuID int64) (*models.Page, error) {
	_ = ctx
	_ = pageID
	_ = skuID
	return nil, errors.New("page not found")
}

func (r *NoopPageDetailsRepository) GetPageImages(ctx context.Context, pageID int64, limit int) ([]models.PageImage, error) {
	_ = ctx
	_ = pageID
	_ = limit
	return nil, nil
}

func (r *NoopPageDetailsRepository) GetBrokenTargetsFrom(ctx context.Context, pageID int64, limit int) ([]models.Page, error) {
	_ = ctx
	_ = pageID
	_ = limit
	return nil, nil
}

func (r *NoopPageDetailsRepository) GetReferrersToBroken(ctx context.Context, pageID int64, limit int) ([]models.Page, error) {
	_ = ctx
	_ = pageID
	_ = limit
	return nil, nil
}

type NoopStatsRepository struct{}

func NewNoopStatsRepository() *NoopStatsRepository {
	return &NoopStatsRepository{}
}

func (r *NoopStatsRepository) Fetch(ctx context.Context, params StatsQueryParams) (map[string]any, error) {
	_ = ctx
	_ = params
	return map[string]any{}, nil
}
