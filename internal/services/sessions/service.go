package sessions

import (
	"context"
	"sitecrawler/newgo/dto"

	sessionsDto "sitecrawler/newgo/dto/sessions"
)

// Service defines all crawling session operations.
type Service interface {
	Create(ctx context.Context, req sessionsDto.CreateCrawlingSessionRequest) (*dto.Response[sessionsDto.CrawlingSessionResponse], error)
	Get(ctx context.Context, req sessionsDto.GetCrawlingSessionRequest) (*dto.Response[sessionsDto.CrawlingSessionResponse], error)
	ListPages(ctx context.Context, req sessionsDto.ListCrawlingSessionPagesRequest) (*dto.Response[sessionsDto.CrawlingSessionPagesResponse], error)
	ListChecks(ctx context.Context, req sessionsDto.ListCrawlingSessionChecksRequest) (*dto.Response[sessionsDto.CrawlingSessionChecksResponse], error)
}
