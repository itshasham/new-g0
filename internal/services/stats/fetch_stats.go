package stats

import (
	"context"
	"net/http"
	"sitecrawler/newgo/dto"

	statsDto "sitecrawler/newgo/dto/stats"
	"sitecrawler/newgo/internal/repository"
)

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
