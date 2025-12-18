package sessions

import (
	"context"
	"net/http"
	"sitecrawler/newgo/dto"

	sessionsDto "sitecrawler/newgo/dto/sessions"
	"sitecrawler/newgo/internal/repository"
	"sitecrawler/newgo/models"
)

type service struct {
	sessionRepo repository.CrawlingSessionRepository
	pageRepo    repository.CrawlingSessionPageRepository
	checkRepo   repository.CrawlingSessionCheckRepository
}

// NewService creates a new crawling session service.
func NewService(
	sessionRepo repository.CrawlingSessionRepository,
	pageRepo repository.CrawlingSessionPageRepository,
	checkRepo repository.CrawlingSessionCheckRepository,
) Service {
	if sessionRepo == nil {
		panic("crawling session repository required")
	}
	if pageRepo == nil {
		panic("page repository required")
	}
	if checkRepo == nil {
		panic("check repository required")
	}
	return &service{
		sessionRepo: sessionRepo,
		pageRepo:    pageRepo,
		checkRepo:   checkRepo,
	}
}

func (s *service) Create(ctx context.Context, req sessionsDto.CreateCrawlingSessionRequest) (*dto.Response[sessionsDto.CrawlingSessionResponse], error) {
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
