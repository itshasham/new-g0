package stats

import (
"sitecrawler/newgo/dto"
	"context"
	"net/http"

	statsDto "sitecrawler/newgo/dto/stats"
	"sitecrawler/newgo/internal/repository"
)

type service struct {
	statsRepo       repository.StatsRepository
	pageDetailsRepo repository.PageDetailsRepository
}

// NewService creates a new stats service.
func NewService(statsRepo repository.StatsRepository, pageDetailsRepo repository.PageDetailsRepository) Service {
	if statsRepo == nil {
		panic("stats repository required")
	}
	if pageDetailsRepo == nil {
		panic("page details repository required")
	}
	return &service{
		statsRepo:       statsRepo,
		pageDetailsRepo: pageDetailsRepo,
	}
}

func (s *service) Fetch(ctx context.Context, req statsDto.StatsRequest) (*dto.Response[statsDto.StatsResponse], error) {
	params := repository.StatsQueryParams{
		SessionID:           req.CrawlingSessionID,
		ComparisonSessionID: req.ComparisonCrawlingID,
		Filters:             req.Filters,
		Prefilters:          req.Prefilters,
	}

	data, err := s.statsRepo.Fetch(ctx, params)
	if err != nil {
		return dto.NewResponse[statsDto.StatsResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	return dto.NewSuccessResponse(statsDto.StatsResponse{Data: data}, http.StatusOK), nil
}
