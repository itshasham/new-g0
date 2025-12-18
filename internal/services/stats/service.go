package stats

import (
	"context"
	"sitecrawler/newgo/controllers/dto"
	"sitecrawler/newgo/internal/repository"

	statsDto "sitecrawler/newgo/controllers/dto/stats"
)

// Service defines all stats operations.
type Service interface {
	Fetch(ctx context.Context, req statsDto.StatsRequest) (*dto.Response[statsDto.StatsResponse], error)
	Details(ctx context.Context, req statsDto.PageDetailsRequest) (*dto.Response[statsDto.PageDetailsResponse], error)
}

type Client struct {
	statsRepo       repository.StatsRepository
	pageDetailsRepo repository.PageDetailsRepository
}

type Option func(s *Client)

// NewService creates a new stats service.
func NewService(opts ...Option) Service {
	client := &Client{}
	client.WithOptions(opts...)
	if client.statsRepo == nil {
		panic("stats repository required")
	}
	if client.pageDetailsRepo == nil {
		panic("page details repository required")
	}
	return client
}

func (c *Client) WithOptions(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithStatsRepository(repo repository.StatsRepository) Option {
	return func(c *Client) {
		c.statsRepo = repo
	}
}

func WithPageDetailsRepository(repo repository.PageDetailsRepository) Option {
	return func(c *Client) {
		c.pageDetailsRepo = repo
	}
}
