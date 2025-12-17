package services

import (
	"context"
	"net/http"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"
)

type viewPageCountService struct {
	viewRepo repository.ViewRepository
	pageRepo repository.CrawlingSessionPageRepository
}

func NewViewPageCountService(viewRepo repository.ViewRepository, pageRepo repository.CrawlingSessionPageRepository) ViewPageCountService {
	if viewRepo == nil {
		panic("view repository required")
	}
	if pageRepo == nil {
		panic("page repository required")
	}
	return &viewPageCountService{viewRepo: viewRepo, pageRepo: pageRepo}
}

func (s *viewPageCountService) PageCount(ctx context.Context, req dto.ViewPageCountRequest) (*dto.Response[dto.ViewPageCountResponse], error) {
	v, err := s.viewRepo.Get(ctx, req.ViewID)
	if err != nil {
		return dto.NewResponse[dto.ViewPageCountResponse](false, err.Error(), http.StatusNotFound, nil), nil
	}

	// Build filter from view's filter_config and use the page repo's List method
	var filters []map[string]any
	if fg, ok := v.FilterConfig["filter_groups"].([]any); ok {
		for _, f := range fg {
			if m, ok := f.(map[string]any); ok {
				filters = append(filters, m)
			}
		}
	}

	params := repository.PageListParams{
		SessionID: req.SessionID,
		Filters:   filters,
		PageLimit: 1,
		Page:      1,
	}

	_, total, err := s.pageRepo.List(ctx, params)
	if err != nil {
		return dto.NewResponse[dto.ViewPageCountResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	return dto.NewSuccessResponse(dto.ViewPageCountResponse{Data: dto.ViewPageCountData{PageCount: total}}, http.StatusOK), nil
}
