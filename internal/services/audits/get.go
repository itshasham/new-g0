package audits

import (
"sitecrawler/newgo/dto"
	"context"
	"errors"
	"net/http"

	auditsDto "sitecrawler/newgo/dto/audits"
	"sitecrawler/newgo/internal/repository"
)

func (s *service) Get(ctx context.Context, req auditsDto.GetAuditCheckRequest) (*dto.Response[auditsDto.AuditCheckResponse], error) {
	check, err := s.repo.Get(ctx, req.ID)
	if err != nil {
		if errors.Is(err, repository.ErrAuditCheckNotFound) {
			return dto.NewResponse[auditsDto.AuditCheckResponse](false, err.Error(), http.StatusNotFound, nil), nil
		}
		return dto.NewResponse[auditsDto.AuditCheckResponse](false, err.Error(), http.StatusInternalServerError, nil), nil
	}

	return dto.NewSuccessResponse(auditsDto.AuditCheckResponse{Data: *check}, http.StatusOK), nil
}
