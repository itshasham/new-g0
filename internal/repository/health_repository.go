package repository

import "context"

// HealthRepository handles database or infrastructure I/O for the health endpoints.
type HealthRepository interface {
	Ping(ctx context.Context) error
}

// NoopHealthRepository implements HealthRepository without hitting external systems.
type NoopHealthRepository struct{}

func NewNoopHealthRepository() *NoopHealthRepository {
	return &NoopHealthRepository{}
}

func (r *NoopHealthRepository) Ping(ctx context.Context) error {
	_ = ctx
	return nil
}
