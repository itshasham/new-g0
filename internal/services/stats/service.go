package stats

import (
	"context"
	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"

	statsDto "sitecrawler/newgo/dto/stats"
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

// Service defines all stats operations.
type Service interface {
	Fetch(ctx context.Context, req statsDto.StatsRequest) (*dto.Response[statsDto.StatsResponse], error)
	Details(ctx context.Context, req statsDto.PageDetailsRequest) (*dto.Response[statsDto.PageDetailsResponse], error)
}
