package services

import (
	"context"
	"net/http"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"
)

type auditCheckListService struct {
	repo repository.AuditCheckRepository
}

func NewAuditCheckListService(repo repository.AuditCheckRepository) AuditCheckListService {
	if repo == nil {
		panic("audit check repository required")
	}
	return &auditCheckListService{repo: repo}
}

func (s *auditCheckListService) List(ctx context.Context, req dto.ListAuditChecksRequest) (*dto.Response[dto.AuditChecksResponse], error) {
	items, err := s.repo.ListBySKU(ctx, req.SearchKeywordURLID)
	if err != nil {
		return dto.NewResponse[dto.AuditChecksResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}
	return dto.NewSuccessResponse(dto.AuditChecksResponse{Data: items}, http.StatusOK), nil
}
