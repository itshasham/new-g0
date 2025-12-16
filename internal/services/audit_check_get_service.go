package services

import (
	"context"
	"errors"
	"net/http"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"
)

type auditCheckGetService struct {
	repo repository.AuditCheckRepository
}

func NewAuditCheckGetService(repo repository.AuditCheckRepository) AuditCheckGetService {
	if repo == nil {
		panic("audit check repository required")
	}
	return &auditCheckGetService{repo: repo}
}

func (s *auditCheckGetService) Get(ctx context.Context, req dto.GetAuditCheckRequest) (*dto.Response[dto.AuditCheckResponse], error) {
	ac, err := s.repo.Get(ctx, req.ID)
	if err != nil {
		if errors.Is(err, repository.ErrAuditCheckNotFound) {
			return dto.NewResponse[dto.AuditCheckResponse](false, err.Error(), http.StatusNotFound, nil), nil
		}
		return nil, err
	}
	return dto.NewSuccessResponse(dto.AuditCheckResponse{Data: *ac}, http.StatusOK), nil
}
