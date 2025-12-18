package audits

import (
"sitecrawler/newgo/dto"
	"context"
	"net/http"

	auditsDto "sitecrawler/newgo/dto/audits"
	"sitecrawler/newgo/models"
)

func (s *service) Create(ctx context.Context, req auditsDto.CreateAuditCheckRequest) (*dto.Response[auditsDto.AuditCheckResponse], error) {
	check := &models.AuditCheck{
		SearchKeywordURLID: req.Data.SearchKeywordURLID,
		Name:               req.Data.Name,
	}

	if err := s.repo.Create(ctx, check); err != nil {
		return dto.NewResponse[auditsDto.AuditCheckResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	return dto.NewSuccessResponse(auditsDto.AuditCheckResponse{Data: *check}, http.StatusCreated), nil
}
