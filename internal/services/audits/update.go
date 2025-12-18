package audits

import (
	"context"
	"errors"
	"net/http"
	"sitecrawler/newgo/dto"

	auditsDto "sitecrawler/newgo/dto/audits"
	"sitecrawler/newgo/internal/repository"
)

func (s *service) Update(ctx context.Context, req auditsDto.UpdateAuditCheckRequest) (*dto.Response[auditsDto.AuditCheckResponse], error) {
	existing, err := s.repo.Get(ctx, req.ID)
	if err != nil {
		if errors.Is(err, repository.ErrAuditCheckNotFound) {
			return dto.NewResponse[auditsDto.AuditCheckResponse](false, err.Error(), http.StatusNotFound, nil), nil
		}
		return dto.NewResponse[auditsDto.AuditCheckResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	if req.Data.Name != nil {
		existing.Name = *req.Data.Name
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		return dto.NewResponse[auditsDto.AuditCheckResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	return dto.NewSuccessResponse(auditsDto.AuditCheckResponse{Data: *existing}, http.StatusOK), nil
}
