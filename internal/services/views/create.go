package views

import (
	"context"
	"net/http"
	"sitecrawler/newgo/controllers/dto"

	viewsDto "sitecrawler/newgo/controllers/dto/views"
	"sitecrawler/newgo/models"
)

func (s *Client) Create(ctx context.Context, req viewsDto.CreateViewRequest) (*dto.Response[viewsDto.ViewResponse], error) {
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
