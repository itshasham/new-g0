package services

import (
	"context"
	"net/http"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"
)

type crawlingSessionPagesService struct {
	repo repository.CrawlingSessionPageRepository
}

func NewCrawlingSessionPagesService(repo repository.CrawlingSessionPageRepository) CrawlingSessionPagesService {
	if repo == nil {
		panic("crawling session page repository required")
	}
	return &crawlingSessionPagesService{repo: repo}
}

func (s *crawlingSessionPagesService) List(ctx context.Context, req dto.ListCrawlingSessionPagesRequest) (*dto.Response[dto.CrawlingSessionPagesResponse], error) {
	pages, total, err := s.repo.List(ctx, repository.PageListParams{
		SessionID: req.SessionID,
		Filters:   req.Filters,
		Sort:      req.Sort,
		Direction: req.Direction,
		Page:      req.Page,
		PageLimit: req.PageLimit,
	})
	if err != nil {
		return dto.NewResponse[dto.CrawlingSessionPagesResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	res := dto.CrawlingSessionPagesResponse{
		Data: dto.CrawlingSessionPagesData{
			Pages:      pages,
			PagesTotal: total,
		},
	}
	return dto.NewSuccessResponse(res, http.StatusOK), nil
}
