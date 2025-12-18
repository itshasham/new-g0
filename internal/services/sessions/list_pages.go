package sessions

import (
	"context"
	"net/http"
	"sitecrawler/newgo/controllers/dto"

	sessionsDto "sitecrawler/newgo/controllers/dto/sessions"
	"sitecrawler/newgo/internal/repository"
)

func (s *Client) ListPages(ctx context.Context, req sessionsDto.ListCrawlingSessionPagesRequest) (*dto.Response[sessionsDto.CrawlingSessionPagesResponse], error) {
	params := repository.PageListParams{
		SessionID: req.SessionID,
		Filters:   req.Filters,
		Sort:      req.Sort,
		Page:      req.Page,
		PageLimit: req.PageLimit,
	}

	pages, total, err := s.pageRepo.List(ctx, params)
	if err != nil {
		return dto.NewResponse[sessionsDto.CrawlingSessionPagesResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	return dto.NewSuccessResponse(sessionsDto.CrawlingSessionPagesResponse{Data: sessionsDto.CrawlingSessionPagesData{Pages: pages, PagesTotal: total}}, http.StatusOK), nil
}
