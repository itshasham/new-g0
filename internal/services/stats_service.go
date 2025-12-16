package services

import (
	"context"
	"net/http"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"
)

type statsService struct {
	repo repository.StatsRepository
}

func NewStatsService(repo repository.StatsRepository) StatsService {
	if repo == nil {
		panic("stats repository required")
	}
	return &statsService{repo: repo}
}

func (s *statsService) Fetch(ctx context.Context, req dto.StatsRequest) (*dto.Response[dto.StatsResponse], error) {
	params := repository.StatsQueryParams{
		SessionID:           req.CrawlingSessionID,
		ComparisonSessionID: req.ComparisonCrawlingID,
		Filters:             req.Filters,
		Prefilters:          req.Prefilters,
	}

	data, err := s.repo.Fetch(ctx, params)
	if err != nil {
		return dto.NewResponse[dto.StatsResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	return dto.NewSuccessResponse(dto.StatsResponse{Data: data}, http.StatusOK), nil
}
