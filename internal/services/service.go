package services

import (
	"context"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"
)

// Services aggregates all domain services so callers can configure dependencies once.
type Services interface {
	Health() HealthService
	CrawlingSessionCreator() CrawlingSessionCreateService
	CrawlingSessionGetter() CrawlingSessionGetService
	CrawlingSessionPages() CrawlingSessionPagesService
	CrawlingSessionChecks() CrawlingSessionChecksService
	PageDetails() PageDetailsService
	Stats() StatsService
}

type HealthService interface {
	Status(ctx context.Context) (dto.HealthResponse, error)
}

type CrawlingSessionCreateService interface {
	Create(ctx context.Context, req dto.CreateCrawlingSessionRequest) (*dto.Response[dto.CrawlingSessionResponse], error)
}

type CrawlingSessionGetService interface {
	Get(ctx context.Context, req dto.GetCrawlingSessionRequest) (*dto.Response[dto.CrawlingSessionResponse], error)
}

type CrawlingSessionPagesService interface {
	List(ctx context.Context, req dto.ListCrawlingSessionPagesRequest) (*dto.Response[dto.CrawlingSessionPagesResponse], error)
}

type CrawlingSessionChecksService interface {
	List(ctx context.Context, req dto.ListCrawlingSessionChecksRequest) (*dto.Response[dto.CrawlingSessionChecksResponse], error)
}

type PageDetailsService interface {
	Details(ctx context.Context, req dto.PageDetailsRequest) (*dto.Response[dto.PageDetailsResponse], error)
}

type StatsService interface {
	Fetch(ctx context.Context, req dto.StatsRequest) (*dto.Response[dto.StatsResponse], error)
}

type service struct {
	health          HealthService
	crawlingCreator CrawlingSessionCreateService
	crawlingGetter  CrawlingSessionGetService
	crawlingPages   CrawlingSessionPagesService
	crawlingChecks  CrawlingSessionChecksService
	pageDetails     PageDetailsService
	statsService    StatsService

	healthRepo          repository.HealthRepository
	crawlingSessionRepo repository.CrawlingSessionRepository
	pageRepo            repository.CrawlingSessionPageRepository
	checkRepo           repository.CrawlingSessionCheckRepository
	pageDetailsRepo     repository.PageDetailsRepository
	statsRepo           repository.StatsRepository
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
	if s.pageRepo == nil {
		s.pageRepo = repository.NewNoopCrawlingSessionPageRepository()
	}
	if s.checkRepo == nil {
		s.checkRepo = repository.NewNoopCrawlingSessionCheckRepository()
	}
	if s.pageDetailsRepo == nil {
		s.pageDetailsRepo = repository.NewNoopPageDetailsRepository()
	}
	if s.statsRepo == nil {
		s.statsRepo = repository.NewNoopStatsRepository()
	}

	s.health = NewHealthService(s.healthRepo)
	s.crawlingCreator = NewCrawlingSessionCreateService(s.crawlingSessionRepo)
	s.crawlingGetter = NewCrawlingSessionGetService(s.crawlingSessionRepo)
	s.crawlingPages = NewCrawlingSessionPagesService(s.pageRepo)
	s.crawlingChecks = NewCrawlingSessionChecksService(s.checkRepo)
	s.pageDetails = NewPageDetailsService(s.pageDetailsRepo)
	s.statsService = NewStatsService(s.statsRepo)
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

func WithCrawlingSessionPageRepository(repo repository.CrawlingSessionPageRepository) Option {
	return func(s *service) {
		s.pageRepo = repo
	}
}

func WithCrawlingSessionCheckRepository(repo repository.CrawlingSessionCheckRepository) Option {
	return func(s *service) {
		s.checkRepo = repo
	}
}

func WithPageDetailsRepository(repo repository.PageDetailsRepository) Option {
	return func(s *service) {
		s.pageDetailsRepo = repo
	}
}

func WithStatsRepository(repo repository.StatsRepository) Option {
	return func(s *service) {
		s.statsRepo = repo
	}
}

func (s *service) Health() HealthService {
	return s.health
}

func (s *service) CrawlingSessionCreator() CrawlingSessionCreateService {
	return s.crawlingCreator
}

func (s *service) CrawlingSessionGetter() CrawlingSessionGetService {
	return s.crawlingGetter
}

func (s *service) CrawlingSessionPages() CrawlingSessionPagesService {
	return s.crawlingPages
}

func (s *service) CrawlingSessionChecks() CrawlingSessionChecksService {
	return s.crawlingChecks
}

func (s *service) PageDetails() PageDetailsService {
	return s.pageDetails
}

func (s *service) Stats() StatsService {
	return s.statsService
}
