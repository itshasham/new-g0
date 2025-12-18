package sessions

import (
	"context"
	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"

	sessionsDto "sitecrawler/newgo/dto/sessions"
)

type service struct {
	sessionRepo repository.CrawlingSessionRepository
	pageRepo    repository.CrawlingSessionPageRepository
	checkRepo   repository.CrawlingSessionCheckRepository
}

// NewService creates a new crawling session service.
func NewService(
	sessionRepo repository.CrawlingSessionRepository,
	pageRepo repository.CrawlingSessionPageRepository,
	checkRepo repository.CrawlingSessionCheckRepository,
) Service {
	if sessionRepo == nil {
		panic("crawling session repository required")
	}
	if pageRepo == nil {
		panic("page repository required")
	}
	if checkRepo == nil {
		panic("check repository required")
	}
	return &service{
		sessionRepo: sessionRepo,
		pageRepo:    pageRepo,
		checkRepo:   checkRepo,
	}
}

// Service defines all crawling session operations.
type Service interface {
	Create(ctx context.Context, req sessionsDto.CreateCrawlingSessionRequest) (*dto.Response[sessionsDto.CrawlingSessionResponse], error)
	Get(ctx context.Context, req sessionsDto.GetCrawlingSessionRequest) (*dto.Response[sessionsDto.CrawlingSessionResponse], error)
	ListPages(ctx context.Context, req sessionsDto.ListCrawlingSessionPagesRequest) (*dto.Response[sessionsDto.CrawlingSessionPagesResponse], error)
	ListChecks(ctx context.Context, req sessionsDto.ListCrawlingSessionChecksRequest) (*dto.Response[sessionsDto.CrawlingSessionChecksResponse], error)
}
