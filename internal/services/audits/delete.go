package audits

import (
	"context"
	"errors"
	"net/http"
	"sitecrawler/newgo/controllers/dto"

	auditsDto "sitecrawler/newgo/controllers/dto/audits"
	"sitecrawler/newgo/internal/repository"
)

func (s *Client) Delete(ctx context.Context, req auditsDto.DeleteAuditCheckRequest) (*dto.Response[auditsDto.DeleteAuditCheckResponse], error) {
	if err := s.repo.Delete(ctx, req.ID); err != nil {
		if errors.Is(err, repository.ErrAuditCheckNotFound) {
			return dto.NewResponse[auditsDto.DeleteAuditCheckResponse](false, err.Error(), http.StatusNotFound, nil), nil
		}
		return dto.NewResponse[auditsDto.DeleteAuditCheckResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	result := auditsDto.DeleteAuditCheckData(req)
	return dto.NewSuccessResponse(auditsDto.DeleteAuditCheckResponse{Data: result}, http.StatusOK), nil
}
