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

// CrawlingSessionRepository describes the data access needed by the service layer.
type CrawlingSessionRepository interface {
	PreventInProgress(ctx context.Context, skuID int64) error
	Create(ctx context.Context, session *models.CrawlingSession) error
	GetByID(ctx context.Context, id int64) (*models.CrawlingSession, error)
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
