package repository

import (
	"context"
	"errors"
	"sync"

	"sitecrawler/newgo/models"
)

var ErrAuditCheckNotFound = errors.New("audit check not found")

type AuditCheckRepository interface {
	Create(ctx context.Context, ac *models.AuditCheck) error
	Update(ctx context.Context, ac *models.AuditCheck) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*models.AuditCheck, error)
	ListBySKU(ctx context.Context, skuID int64) ([]models.AuditCheck, error)
	ListByIDsAndSKUs(ctx context.Context, ids, skuIDs []int64) ([]models.AuditCheck, error)
}

type InMemoryAuditCheckRepository struct {
	mu    sync.Mutex
	seq   int64
	items map[int64]*models.AuditCheck
}

func NewInMemoryAuditCheckRepository() *InMemoryAuditCheckRepository {
	return &InMemoryAuditCheckRepository{items: map[int64]*models.AuditCheck{}}
}

func (r *InMemoryAuditCheckRepository) Create(ctx context.Context, ac *models.AuditCheck) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	ac.ID = r.seq
	r.items[ac.ID] = cloneAuditCheck(ac)
	return nil
}

func (r *InMemoryAuditCheckRepository) Update(ctx context.Context, ac *models.AuditCheck) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[ac.ID] = cloneAuditCheck(ac)
	return nil
}

func (r *InMemoryAuditCheckRepository) Delete(ctx context.Context, id int64) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}

func (r *InMemoryAuditCheckRepository) Get(ctx context.Context, id int64) (*models.AuditCheck, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	if v, ok := r.items[id]; ok {
		return cloneAuditCheck(v), nil
	}
	return nil, ErrAuditCheckNotFound
}

func (r *InMemoryAuditCheckRepository) ListBySKU(ctx context.Context, skuID int64) ([]models.AuditCheck, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []models.AuditCheck{}
	for _, v := range r.items {
		if v.SearchKeywordURLID == skuID {
			out = append(out, *cloneAuditCheck(v))
		}
	}
	return out, nil
}

func (r *InMemoryAuditCheckRepository) ListByIDsAndSKUs(ctx context.Context, ids, skuIDs []int64) ([]models.AuditCheck, error) {
	_ = ctx
	_ = ids
	_ = skuIDs
	return []models.AuditCheck{}, nil
}

func cloneAuditCheck(ac *models.AuditCheck) *models.AuditCheck {
	c := *ac
	return &c
}
