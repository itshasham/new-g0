package views

import (
	"context"
	"sitecrawler/newgo/controllers/dto"
	"sitecrawler/newgo/internal/repository"

	viewsDto "sitecrawler/newgo/controllers/dto/views"
)

// Service defines all view operations.
type Service interface {
	List(ctx context.Context, req viewsDto.ListViewsRequest) (*dto.Response[viewsDto.ViewsResponse], error)
	Create(ctx context.Context, req viewsDto.CreateViewRequest) (*dto.Response[viewsDto.ViewResponse], error)
	Get(ctx context.Context, req viewsDto.GetViewRequest) (*dto.Response[viewsDto.ViewResponse], error)
	Update(ctx context.Context, req viewsDto.UpdateViewRequest) (*dto.Response[viewsDto.ViewResponse], error)
	Delete(ctx context.Context, req viewsDto.DeleteViewRequest) (*dto.Response[viewsDto.DeleteViewResponse], error)
	PageCount(ctx context.Context, req viewsDto.ViewPageCountRequest) (*dto.Response[viewsDto.ViewPageCountResponse], error)
}

type Client struct {
	viewRepo repository.ViewRepository
	pageRepo repository.CrawlingSessionPageRepository
}

type Option func(s *Client)

// NewService creates a new view service.
func NewService(opts ...Option) Service {
	client := &Client{}
	client.WithOptions(opts...)
	if client.viewRepo == nil {
		panic("view repository required")
	}
	if client.pageRepo == nil {
		panic("page repository required")
	}
	return client
}

func (c *Client) WithOptions(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithViewRepository(repo repository.ViewRepository) Option {
	return func(c *Client) {
		c.viewRepo = repo
	}
}

func WithPageRepository(repo repository.CrawlingSessionPageRepository) Option {
	return func(c *Client) {
		c.pageRepo = repo
	}
}
