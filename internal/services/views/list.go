package views

import (
	"context"
	"net/http"
	"sitecrawler/newgo/dto"

	viewsDto "sitecrawler/newgo/dto/views"
	"sitecrawler/newgo/internal/repository"
	"sitecrawler/newgo/models"
)

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
