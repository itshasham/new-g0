package services

import (
	"context"
	"errors"
	"net/http"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"
)

type crawlingSessionGetService struct {
	repo repository.CrawlingSessionRepository
}

func NewCrawlingSessionGetService(repo repository.CrawlingSessionRepository) CrawlingSessionGetService {
	if repo == nil {
		panic("crawling session repository required")
	}
	return &crawlingSessionGetService{repo: repo}
}

func (s *crawlingSessionGetService) Get(ctx context.Context, req dto.GetCrawlingSessionRequest) (*dto.Response[dto.CrawlingSessionResponse], error) {
	session, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, repository.ErrCrawlingSessionNotFound) {
			return dto.NewResponse[dto.CrawlingSessionResponse](false, err.Error(), http.StatusNotFound, nil), nil
		}
		return nil, err
	}

	return dto.NewSuccessResponse(dto.CrawlingSessionResponse{Data: *session}, http.StatusOK), nil
}
