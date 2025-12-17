package services

import (
	"context"
	"net/http"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"
	"sitecrawler/newgo/models"
)

type viewCreateService struct {
	repo repository.ViewRepository
}

func NewViewCreateService(repo repository.ViewRepository) ViewCreateService {
	if repo == nil {
		panic("view repository required")
	}
	return &viewCreateService{repo: repo}
}

func (s *viewCreateService) Create(ctx context.Context, req dto.CreateViewRequest) (*dto.Response[dto.ViewResponse], error) {
	v := &models.View{
		SearchKeywordURLID: req.Data.SearchKeywordURLID,
		Name:               req.Data.Name,
		FilterConfig:       req.Data.FilterConfig,
	}
	if err := s.repo.Create(ctx, v); err != nil {
		return dto.NewResponse[dto.ViewResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}
	return dto.NewSuccessResponse(dto.ViewResponse{Data: *v}, http.StatusCreated), nil
}
