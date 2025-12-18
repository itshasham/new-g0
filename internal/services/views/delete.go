package views

import (
	"context"
	"errors"
	"net/http"
	"sitecrawler/newgo/controllers/dto"

	viewsDto "sitecrawler/newgo/controllers/dto/views"
	"sitecrawler/newgo/internal/repository"
)

func (s *Client) Delete(ctx context.Context, req viewsDto.DeleteViewRequest) (*dto.Response[viewsDto.DeleteViewResponse], error) {
	if err := s.viewRepo.Delete(ctx, req.ID); err != nil {
		if errors.Is(err, repository.ErrViewNotFound) {
			return dto.NewResponse[viewsDto.DeleteViewResponse](false, err.Error(), http.StatusNotFound, nil), nil
		}
		return dto.NewResponse[viewsDto.DeleteViewResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	result := viewsDto.DeleteViewData{ID: req.ID}
	return dto.NewSuccessResponse(viewsDto.DeleteViewResponse{Data: result}, http.StatusOK), nil
}
