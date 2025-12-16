package services

import (
	"context"
	"net/http"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"
	"sitecrawler/newgo/models"
)

type auditCheckCreateService struct {
	repo repository.AuditCheckRepository
}

func NewAuditCheckCreateService(repo repository.AuditCheckRepository) AuditCheckCreateService {
	if repo == nil {
		panic("audit check repository required")
	}
	return &auditCheckCreateService{repo: repo}
}

func (s *auditCheckCreateService) Create(ctx context.Context, req dto.CreateAuditCheckRequest) (*dto.Response[dto.AuditCheckResponse], error) {
	ac := &models.AuditCheck{
		SearchKeywordURLID: req.Data.SearchKeywordURLID,
		Name:               req.Data.Name,
		Category:           req.Data.Category,
		FilterConfig:       req.Data.FilterConfig,
	}
	if err := s.repo.Create(ctx, ac); err != nil {
		return dto.NewResponse[dto.AuditCheckResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}
	return dto.NewSuccessResponse(dto.AuditCheckResponse{Data: *ac}, http.StatusCreated), nil
}
