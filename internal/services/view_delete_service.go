package services

import (
	"context"
	"net/http"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"
)

type viewDeleteService struct {
	repo repository.ViewRepository
}

func NewViewDeleteService(repo repository.ViewRepository) ViewDeleteService {
	if repo == nil {
		panic("view repository required")
	}
	return &viewDeleteService{repo: repo}
}

func (s *viewDeleteService) Delete(ctx context.Context, req dto.DeleteViewRequest) (*dto.Response[dto.DeleteViewResponse], error) {
	if err := s.repo.Delete(ctx, req.ID); err != nil {
		return dto.NewResponse[dto.DeleteViewResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}
	return dto.NewSuccessResponse(dto.DeleteViewResponse{Data: dto.DeleteViewData{ID: req.ID}}, http.StatusOK), nil
}
