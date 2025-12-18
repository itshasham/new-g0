package sessions

import (
	"context"
	"sitecrawler/newgo/controllers/dto"
	"sitecrawler/newgo/internal/repository"

	sessionsDto "sitecrawler/newgo/controllers/dto/sessions"
)

// Service defines all crawling session operations.
type Service interface {
	Create(ctx context.Context, req sessionsDto.CreateCrawlingSessionRequest) (*dto.Response[sessionsDto.CrawlingSessionResponse], error)
	Get(ctx context.Context, req sessionsDto.GetCrawlingSessionRequest) (*dto.Response[sessionsDto.CrawlingSessionResponse], error)
	ListPages(ctx context.Context, req sessionsDto.ListCrawlingSessionPagesRequest) (*dto.Response[sessionsDto.CrawlingSessionPagesResponse], error)
	ListChecks(ctx context.Context, req sessionsDto.ListCrawlingSessionChecksRequest) (*dto.Response[sessionsDto.CrawlingSessionChecksResponse], error)
}

type Client struct {
	sessionRepo repository.CrawlingSessionRepository
	pageRepo    repository.CrawlingSessionPageRepository
	checkRepo   repository.CrawlingSessionCheckRepository
}

type Option func(s *Client)

// NewService creates a new crawling session service.
func NewService(opts ...Option) Service {
	client := &Client{}
	client.WithOptions(opts...)
	if client.sessionRepo == nil {
		panic("crawling session repository required")
	}
	if client.pageRepo == nil {
		panic("page repository required")
	}
	if client.checkRepo == nil {
		panic("check repository required")
	}
	return client
}

func (c *Client) WithOptions(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithSessionRepository(repo repository.CrawlingSessionRepository) Option {
	return func(c *Client) {
		c.sessionRepo = repo
	}
}

func WithPageRepository(repo repository.CrawlingSessionPageRepository) Option {
	return func(c *Client) {
		c.pageRepo = repo
	}
}

func WithCheckRepository(repo repository.CrawlingSessionCheckRepository) Option {
	return func(c *Client) {
		c.checkRepo = repo
	}
}
