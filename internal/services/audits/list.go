package audits

import (
"sitecrawler/newgo/dto"
	"context"
	"net/http"

	auditsDto "sitecrawler/newgo/dto/audits"
	"sitecrawler/newgo/internal/repository"
)

type service struct {
	repo repository.AuditCheckRepository
}

// NewService creates a new audit check service.
func NewService(repo repository.AuditCheckRepository) Service {
	if repo == nil {
		panic("audit check repository required")
	}
	return &service{repo: repo}
}

func (s *service) List(ctx context.Context, req auditsDto.ListAuditChecksRequest) (*dto.Response[auditsDto.AuditChecksResponse], error) {
	checks, err := s.repo.ListBySKU(ctx, req.SearchKeywordURLID)
	if err != nil {
		return dto.NewResponse[auditsDto.AuditChecksResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	return dto.NewSuccessResponse(auditsDto.AuditChecksResponse{Data: checks}, http.StatusOK), nil
}
