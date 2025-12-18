package views

import (
	"context"
	"net/http"
	"sitecrawler/newgo/controllers/dto"

	viewsDto "sitecrawler/newgo/controllers/dto/views"
	"sitecrawler/newgo/internal/repository"
)

func (s *Client) PageCount(ctx context.Context, req viewsDto.ViewPageCountRequest) (*dto.Response[viewsDto.ViewPageCountResponse], error) {
	v, err := s.viewRepo.Get(ctx, req.ViewID)
	if err != nil {
		return dto.NewResponse[viewsDto.ViewPageCountResponse](false, err.Error(), http.StatusNotFound, nil), nil
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
		return dto.NewResponse[viewsDto.ViewPageCountResponse](false, err.Error(), http.StatusUnprocessableEntity, nil), nil
	}

	return dto.NewSuccessResponse(viewsDto.ViewPageCountResponse{Data: viewsDto.ViewPageCountData{PageCount: total}}, http.StatusOK), nil
}
