package services

import (
	"context"
	"net/http"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"
)

type auditCheckDeleteService struct {
	repo repository.AuditCheckRepository
}

func NewAuditCheckDeleteService(repo repository.AuditCheckRepository) AuditCheckDeleteService {
	if repo == nil {
		panic("audit check repository required")
	}
	return &auditCheckDeleteService{repo: repo}
}

func (s *auditCheckDeleteService) Delete(ctx context.Context, req dto.DeleteAuditCheckRequest) (*dto.Response[dto.DeleteAuditCheckResponse], error) {
	if err := s.repo.Delete(ctx, req.ID); err != nil {
		return dto.NewResponse[dto.DeleteAuditCheckResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}
	return dto.NewSuccessResponse(dto.DeleteAuditCheckResponse{Data: dto.DeleteAuditCheckData{ID: req.ID}}, http.StatusOK), nil
}
