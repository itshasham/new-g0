package services

import (
	"context"
	"net/http"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"
	"sitecrawler/newgo/models"
)

type crawlingSessionService struct {
	repo repository.CrawlingSessionRepository
}

func NewCrawlingSessionService(repo repository.CrawlingSessionRepository) CrawlingSessionService {
	if repo == nil {
		panic("crawling session repository required")
	}
	return &crawlingSessionService{repo: repo}
}

func (s *crawlingSessionService) Create(ctx context.Context, req dto.CreateCrawlingSessionRequest) (*dto.Response[dto.CrawlingSessionResponse], error) {
	if err := s.repo.PreventInProgress(ctx, req.Data.SearchKeywordURLID); err != nil {
		return dto.NewResponse[dto.CrawlingSessionResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	session := &models.CrawlingSession{
		SearchKeywordURLID: req.Data.SearchKeywordURLID,
		URL:                req.Data.URL,
		Options:            req.Data.Options,
		Status:             "pending",
		Queue:              req.Data.Queue,
		Version:            1,
	}

	if err := s.repo.Create(ctx, session); err != nil {
		return nil, err
	}

	return dto.NewSuccessResponse(dto.CrawlingSessionResponse{Data: *session}, http.StatusCreated), nil
}
