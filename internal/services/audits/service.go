package audits

import (
	"context"
	"sitecrawler/newgo/controllers/dto"

	auditsDto "sitecrawler/newgo/controllers/dto/audits"
	"sitecrawler/newgo/internal/repository"
)

// Service defines all audit check operations.
type Service interface {
	List(ctx context.Context, req auditsDto.ListAuditChecksRequest) (*dto.Response[auditsDto.AuditChecksResponse], error)
	Create(ctx context.Context, req auditsDto.CreateAuditCheckRequest) (*dto.Response[auditsDto.AuditCheckResponse], error)
	Get(ctx context.Context, req auditsDto.GetAuditCheckRequest) (*dto.Response[auditsDto.AuditCheckResponse], error)
	Update(ctx context.Context, req auditsDto.UpdateAuditCheckRequest) (*dto.Response[auditsDto.AuditCheckResponse], error)
	Delete(ctx context.Context, req auditsDto.DeleteAuditCheckRequest) (*dto.Response[auditsDto.DeleteAuditCheckResponse], error)
}

type Client struct {
	repo repository.AuditCheckRepository
}

type Option func(s *Client)

func NewService(opts ...Option) Service {
	client := &Client{}
	client.WithOptions(opts...)
	if client.repo == nil {
		panic("audit check repository required")
	}
	return client
}

func (c *Client) WithOptions(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithAuditCheckRepository(repo repository.AuditCheckRepository) Option {
	return func(c *Client) {
		c.repo = repo
	}
}
