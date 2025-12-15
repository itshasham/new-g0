package services

import (
	"context"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"
)

// Services aggregates all domain services so callers can configure dependencies once.
type Services interface {
	Health() HealthService
	CrawlingSessions() CrawlingSessionService
}

type HealthService interface {
	Status(ctx context.Context) (dto.HealthResponse, error)
}

type CrawlingSessionService interface {
	Create(ctx context.Context, req dto.CreateCrawlingSessionRequest) (*dto.Response[dto.CrawlingSessionResponse], error)
}

type service struct {
	health           HealthService
	crawlingSessions CrawlingSessionService

	healthRepo          repository.HealthRepository
	crawlingSessionRepo repository.CrawlingSessionRepository
}

// Option allows callers to configure the Services container.
type Option func(*service)

// New constructs the Services implementation using provided options.
func New(opts ...Option) Services {
	s := &service{}
	for _, opt := range opts {
		opt(s)
	}

	if s.healthRepo == nil {
		s.healthRepo = repository.NewNoopHealthRepository()
	}
	if s.crawlingSessionRepo == nil {
		s.crawlingSessionRepo = repository.NewInMemoryCrawlingSessionRepository()
	}

	s.health = NewHealthService(s.healthRepo)
	s.crawlingSessions = NewCrawlingSessionService(s.crawlingSessionRepo)
	return s
}

// WithHealthRepository overrides the default health repository.
func WithHealthRepository(repo repository.HealthRepository) Option {
	return func(s *service) {
		s.healthRepo = repo
	}
}

// WithCrawlingSessionRepository overrides the default crawling session repository.
func WithCrawlingSessionRepository(repo repository.CrawlingSessionRepository) Option {
	return func(s *service) {
		s.crawlingSessionRepo = repo
	}
}

func (s *service) Health() HealthService {
	return s.health
}

func (s *service) CrawlingSessions() CrawlingSessionService {
	return s.crawlingSessions
}
