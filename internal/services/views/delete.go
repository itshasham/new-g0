package views

import (
	"context"
	"errors"
	"net/http"
	"sitecrawler/newgo/dto"

	viewsDto "sitecrawler/newgo/dto/views"
	"sitecrawler/newgo/internal/repository"
)

func (s *service) Delete(ctx context.Context, req viewsDto.DeleteViewRequest) (*dto.Response[viewsDto.DeleteViewResponse], error) {
	if err := s.viewRepo.Delete(ctx, req.ID); err != nil {
		if errors.Is(err, repository.ErrViewNotFound) {
			return dto.NewResponse[viewsDto.DeleteViewResponse](false, err.Error(), http.StatusNotFound, nil), nil
		}
		return dto.NewResponse[viewsDto.DeleteViewResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	result := viewsDto.DeleteViewData{ID: req.ID}
	return dto.NewSuccessResponse(viewsDto.DeleteViewResponse{Data: result}, http.StatusOK), nil
}
