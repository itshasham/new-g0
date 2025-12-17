package services

import (
	"context"
	"net/http"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"
	"sitecrawler/newgo/models"
)

type viewListService struct {
	repo repository.ViewRepository
}

func NewViewListService(repo repository.ViewRepository) ViewListService {
	if repo == nil {
		panic("view repository required")
	}
	return &viewListService{repo: repo}
}

func (s *viewListService) List(ctx context.Context, req dto.ListViewsRequest) (*dto.Response[dto.ViewsResponse], error) {
	if req.SearchKeywordURLID == 0 {
		return dto.NewSuccessResponse(dto.ViewsResponse{Data: []models.View{}}, http.StatusOK), nil
	}
	items, err := s.repo.ListBySKU(ctx, req.SearchKeywordURLID)
	if err != nil {
		return dto.NewResponse[dto.ViewsResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}
	return dto.NewSuccessResponse(dto.ViewsResponse{Data: items}, http.StatusOK), nil
}
