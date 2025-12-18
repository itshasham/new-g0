package audits

import (
	"context"
	"net/http"
	"sitecrawler/newgo/dto"

	auditsDto "sitecrawler/newgo/dto/audits"
)

func (s *service) List(ctx context.Context, req auditsDto.ListAuditChecksRequest) (*dto.Response[auditsDto.AuditChecksResponse], error) {
	checks, err := s.repo.ListBySKU(ctx, req.SearchKeywordURLID)
	if err != nil {
		return dto.NewResponse[auditsDto.AuditChecksResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	return dto.NewSuccessResponse(auditsDto.AuditChecksResponse{Data: checks}, http.StatusOK), nil
}
