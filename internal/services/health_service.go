package services

import (
	"context"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"
)

type healthService struct {
	repo repository.HealthRepository
}

func NewHealthService(repo repository.HealthRepository) HealthService {
	if repo == nil {
		panic("repository must not be nil")
	}

	return &healthService{
		repo: repo,
	}
}

func (s *healthService) Status(ctx context.Context) (dto.HealthResponse, error) {
	if err := s.repo.Ping(ctx); err != nil {
		return dto.HealthResponse{}, err
	}

	return dto.HealthResponse{Status: "ok"}, nil
}
