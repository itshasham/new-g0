package views

import (
"sitecrawler/newgo/dto"
	"context"
	"net/http"

	viewsDto "sitecrawler/newgo/dto/views"
	"sitecrawler/newgo/models"
)

func (s *service) Create(ctx context.Context, req viewsDto.CreateViewRequest) (*dto.Response[viewsDto.ViewResponse], error) {
	view := &models.View{
		SearchKeywordURLID: req.Data.SearchKeywordURLID,
		Name:               req.Data.Name,
		FilterConfig:       req.Data.FilterConfig,
	}

	if err := s.viewRepo.Create(ctx, view); err != nil {
		return dto.NewResponse[viewsDto.ViewResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	return dto.NewSuccessResponse(viewsDto.ViewResponse{Data: *view}, http.StatusCreated), nil
}
