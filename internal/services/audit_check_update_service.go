package services

import (
	"context"
	"errors"
	"net/http"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"
)

type auditCheckUpdateService struct {
	repo repository.AuditCheckRepository
}

func NewAuditCheckUpdateService(repo repository.AuditCheckRepository) AuditCheckUpdateService {
	if repo == nil {
		panic("audit check repository required")
	}
	return &auditCheckUpdateService{repo: repo}
}

func (s *auditCheckUpdateService) Update(ctx context.Context, req dto.UpdateAuditCheckRequest) (*dto.Response[dto.AuditCheckResponse], error) {
	ac, err := s.repo.Get(ctx, req.ID)
	if err != nil {
		if errors.Is(err, repository.ErrAuditCheckNotFound) {
			return dto.NewResponse[dto.AuditCheckResponse](false, err.Error(), http.StatusNotFound, nil), nil
		}
		return nil, err
	}

	if req.Data.Name != nil {
		ac.Name = *req.Data.Name
	}
	if req.Data.Category != nil {
		ac.Category = *req.Data.Category
	}
	if req.Data.FilterConfig != nil {
		ac.FilterConfig = *req.Data.FilterConfig
	}

	if err := s.repo.Update(ctx, ac); err != nil {
		return dto.NewResponse[dto.AuditCheckResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}
	return dto.NewSuccessResponse(dto.AuditCheckResponse{Data: *ac}, http.StatusOK), nil
}
