package views

import (
"sitecrawler/newgo/dto"
	"context"
	"net/http"

	viewsDto "sitecrawler/newgo/dto/views"
	"sitecrawler/newgo/internal/repository"
	"sitecrawler/newgo/models"
)

type service struct {
	viewRepo repository.ViewRepository
	pageRepo repository.CrawlingSessionPageRepository
}

// NewService creates a new view service.
func NewService(viewRepo repository.ViewRepository, pageRepo repository.CrawlingSessionPageRepository) Service {
	if viewRepo == nil {
		panic("view repository required")
	}
	if pageRepo == nil {
		panic("page repository required")
	}
	return &service{
		viewRepo: viewRepo,
		pageRepo: pageRepo,
	}
}

func (s *service) List(ctx context.Context, req viewsDto.ListViewsRequest) (*dto.Response[viewsDto.ViewsResponse], error) {
	skuID := req.SearchKeywordURLID
	if skuID == 0 {
		return dto.NewSuccessResponse(viewsDto.ViewsResponse{Data: []models.View{}}, http.StatusOK), nil
	}

	views, err := s.viewRepo.ListBySKU(ctx, skuID)
	if err != nil {
		return dto.NewResponse[viewsDto.ViewsResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	return dto.NewSuccessResponse(viewsDto.ViewsResponse{Data: views}, http.StatusOK), nil
}
