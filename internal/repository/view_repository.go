package repository

import (
	"context"
	"errors"
	"sync"

	"sitecrawler/newgo/models"
)

var ErrViewNotFound = errors.New("view not found")

type ViewRepository interface {
	Create(ctx context.Context, v *models.View) error
	Update(ctx context.Context, v *models.View) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*models.View, error)
	ListBySKU(ctx context.Context, skuID int64) ([]models.View, error)
}

type InMemoryViewRepository struct {
	mu    sync.Mutex
	seq   int64
	items map[int64]*models.View
}

func NewInMemoryViewRepository() *InMemoryViewRepository {
	return &InMemoryViewRepository{items: map[int64]*models.View{}}
}

func (r *InMemoryViewRepository) Create(ctx context.Context, v *models.View) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	v.ID = r.seq
	r.items[v.ID] = cloneView(v)
	return nil
}

func (r *InMemoryViewRepository) Update(ctx context.Context, v *models.View) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[v.ID]; !ok {
		return ErrViewNotFound
	}
	r.items[v.ID] = cloneView(v)
	return nil
}

func (r *InMemoryViewRepository) Delete(ctx context.Context, id int64) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}

func (r *InMemoryViewRepository) Get(ctx context.Context, id int64) (*models.View, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	if v, ok := r.items[id]; ok {
		return cloneView(v), nil
	}
	return nil, ErrViewNotFound
}

func (r *InMemoryViewRepository) ListBySKU(ctx context.Context, skuID int64) ([]models.View, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []models.View{}
	for _, v := range r.items {
		if v.SearchKeywordURLID == skuID {
			out = append(out, *cloneView(v))
		}
	}
	return out, nil
}

func cloneView(v *models.View) *models.View {
	c := *v
	return &c
}
