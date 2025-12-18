package views

import (
	"context"
	"net/http"
	"sitecrawler/newgo/dto"

	viewsDto "sitecrawler/newgo/dto/views"
)

func (s *service) Get(ctx context.Context, req viewsDto.GetViewRequest) (*dto.Response[viewsDto.ViewResponse], error) {
	view, err := s.viewRepo.Get(ctx, req.ID)
	if err != nil {
		// Return 404 for any error (including repository failures)
		return dto.NewResponse[viewsDto.ViewResponse](false, err.Error(), http.StatusNotFound, nil), nil
	}

	return dto.NewSuccessResponse(viewsDto.ViewResponse{Data: *view}, http.StatusOK), nil
}
