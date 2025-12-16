package services

import (
	"context"
	"net/http"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"
)

type crawlingSessionChecksService struct {
	repo repository.CrawlingSessionCheckRepository
}

func NewCrawlingSessionChecksService(repo repository.CrawlingSessionCheckRepository) CrawlingSessionChecksService {
	if repo == nil {
		panic("crawling session check repository required")
	}
	return &crawlingSessionChecksService{repo: repo}
}

func (s *crawlingSessionChecksService) List(ctx context.Context, req dto.ListCrawlingSessionChecksRequest) (*dto.Response[dto.CrawlingSessionChecksResponse], error) {
	checks, err := s.repo.ChecksWithPages(ctx, repository.ChecksWithPagesParams{
		SessionID:           req.SessionID,
		ComparisonSessionID: req.ComparisonSessionID,
		ViewFilters:         req.ViewFilters,
		PageLimitPerCheck:   req.PageLimitPerCheck,
	})
	if err != nil {
		return dto.NewResponse[dto.CrawlingSessionChecksResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	res := dto.CrawlingSessionChecksResponse{
		Data: dto.CrawlingSessionChecksData{
			Checks: checks,
		},
	}
	return dto.NewSuccessResponse(res, http.StatusOK), nil
}
