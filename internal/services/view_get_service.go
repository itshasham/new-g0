package services

import (
	"context"
	"net/http"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"
)

type viewGetService struct {
	repo repository.ViewRepository
}

func NewViewGetService(repo repository.ViewRepository) ViewGetService {
	if repo == nil {
		panic("view repository required")
	}
	return &viewGetService{repo: repo}
}

func (s *viewGetService) Get(ctx context.Context, req dto.GetViewRequest) (*dto.Response[dto.ViewResponse], error) {
	v, err := s.repo.Get(ctx, req.ID)
	if err != nil {
		return dto.NewResponse[dto.ViewResponse](false, err.Error(), http.StatusNotFound, nil), nil
	}
	return dto.NewSuccessResponse(dto.ViewResponse{Data: *v}, http.StatusOK), nil
}
