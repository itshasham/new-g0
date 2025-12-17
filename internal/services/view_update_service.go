package services

import (
	"context"
	"net/http"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"
)

type viewUpdateService struct {
	repo repository.ViewRepository
}

func NewViewUpdateService(repo repository.ViewRepository) ViewUpdateService {
	if repo == nil {
		panic("view repository required")
	}
	return &viewUpdateService{repo: repo}
}

func (s *viewUpdateService) Update(ctx context.Context, req dto.UpdateViewRequest) (*dto.Response[dto.ViewResponse], error) {
	v, err := s.repo.Get(ctx, req.ID)
	if err != nil {
		return dto.NewResponse[dto.ViewResponse](false, err.Error(), http.StatusNotFound, nil), nil
	}

	if req.Data.Name != nil {
		v.Name = *req.Data.Name
	}
	if req.Data.FilterConfig != nil {
		v.FilterConfig = *req.Data.FilterConfig
	}

	if err := s.repo.Update(ctx, v); err != nil {
		return dto.NewResponse[dto.ViewResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}
	return dto.NewSuccessResponse(dto.ViewResponse{Data: *v}, http.StatusOK), nil
}
