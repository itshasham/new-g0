package sessions

import (
	"context"
	"errors"
	"net/http"
	"sitecrawler/newgo/dto"

	sessionsDto "sitecrawler/newgo/dto/sessions"
	"sitecrawler/newgo/internal/repository"
)

func (s *service) Get(ctx context.Context, req sessionsDto.GetCrawlingSessionRequest) (*dto.Response[sessionsDto.CrawlingSessionResponse], error) {
	session, err := s.sessionRepo.GetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, repository.ErrCrawlingSessionNotFound) {
			return dto.NewResponse[sessionsDto.CrawlingSessionResponse](false, err.Error(), http.StatusNotFound, nil), nil
		}
		return dto.NewResponse[sessionsDto.CrawlingSessionResponse](false, err.Error(), http.StatusInternalServerError, nil), nil
	}

	return dto.NewSuccessResponse(sessionsDto.CrawlingSessionResponse{Data: *session}, http.StatusOK), nil
}
