package sessions

import (
	"context"
	"net/http"
	"sitecrawler/newgo/dto"

	sessionsDto "sitecrawler/newgo/dto/sessions"
	"sitecrawler/newgo/internal/repository"
)

func (s *service) ListChecks(ctx context.Context, req sessionsDto.ListCrawlingSessionChecksRequest) (*dto.Response[sessionsDto.CrawlingSessionChecksResponse], error) {
	params := repository.ChecksWithPagesParams{
		SessionID:           req.SessionID,
		ComparisonSessionID: req.ComparisonSessionID,
		ViewFilters:         req.ViewFilters,
		PageLimitPerCheck:   req.PageLimitPerCheck,
	}

	checks, err := s.checkRepo.ChecksWithPages(ctx, params)
	if err != nil {
		return dto.NewResponse[sessionsDto.CrawlingSessionChecksResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	return dto.NewSuccessResponse(sessionsDto.CrawlingSessionChecksResponse{Data: sessionsDto.CrawlingSessionChecksData{Checks: checks}}, http.StatusOK), nil
}
