package sessions

import (
	"context"
	"net/http"
	"sitecrawler/newgo/controllers/dto"

	sessionsDto "sitecrawler/newgo/controllers/dto/sessions"
	"sitecrawler/newgo/models"
)

func (s *Client) Create(ctx context.Context, req sessionsDto.CreateCrawlingSessionRequest) (*dto.Response[sessionsDto.CrawlingSessionResponse], error) {
	if err := s.sessionRepo.PreventInProgress(ctx, req.Data.SearchKeywordURLID); err != nil {
		return dto.NewResponse[sessionsDto.CrawlingSessionResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	session := &models.CrawlingSession{
		SearchKeywordURLID: req.Data.SearchKeywordURLID,
		Status:             "pending",
		URL:                req.Data.URL,
		Queue:              req.Data.Queue,
		Options:            req.Data.Options,
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return dto.NewResponse[sessionsDto.CrawlingSessionResponse](false, err.Error(), http.StatusInternalServerError, nil), nil
	}

	return dto.NewSuccessResponse(sessionsDto.CrawlingSessionResponse{Data: *session}, http.StatusCreated), nil
}
