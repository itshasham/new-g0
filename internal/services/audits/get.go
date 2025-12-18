package audits

import (
	"context"
	"errors"
	"net/http"
	"sitecrawler/newgo/controllers/dto"

	auditsDto "sitecrawler/newgo/controllers/dto/audits"
	"sitecrawler/newgo/internal/repository"
)

func (s *Client) Get(ctx context.Context, req auditsDto.GetAuditCheckRequest) (*dto.Response[auditsDto.AuditCheckResponse], error) {
	check, err := s.repo.Get(ctx, req.ID)
	if err != nil {
		if errors.Is(err, repository.ErrAuditCheckNotFound) {
			return dto.NewResponse[auditsDto.AuditCheckResponse](false, err.Error(), http.StatusNotFound, nil), nil
		}
		return dto.NewResponse[auditsDto.AuditCheckResponse](false, err.Error(), http.StatusInternalServerError, nil), nil
	}

	return dto.NewSuccessResponse(auditsDto.AuditCheckResponse{Data: *check}, http.StatusOK), nil
}
