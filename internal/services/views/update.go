package views

import (
"sitecrawler/newgo/dto"
	"context"
	"errors"
	"net/http"

	viewsDto "sitecrawler/newgo/dto/views"
	"sitecrawler/newgo/internal/repository"
)

func (s *service) Update(ctx context.Context, req viewsDto.UpdateViewRequest) (*dto.Response[viewsDto.ViewResponse], error) {
	existing, err := s.viewRepo.Get(ctx, req.ID)
	if err != nil {
		if errors.Is(err, repository.ErrViewNotFound) {
			return dto.NewResponse[viewsDto.ViewResponse](false, err.Error(), http.StatusNotFound, nil), nil
		}
		return dto.NewResponse[viewsDto.ViewResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	if req.Data.Name != nil && *req.Data.Name != "" {
		existing.Name = *req.Data.Name
	}
	if req.Data.FilterConfig != nil {
		existing.FilterConfig = *req.Data.FilterConfig
	}

	if err := s.viewRepo.Update(ctx, existing); err != nil {
		return dto.NewResponse[viewsDto.ViewResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	return dto.NewSuccessResponse(viewsDto.ViewResponse{Data: *existing}, http.StatusOK), nil
}
