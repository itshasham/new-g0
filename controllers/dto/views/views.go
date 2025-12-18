package views

import "sitecrawler/newgo/models"

type ListViewsRequest struct {
	SearchKeywordURLID int64 `json:"search_keyword_url_id"`
}

type ViewResponse struct {
	Data models.View `json:"data"`
}

type ViewsResponse struct {
	Data []models.View `json:"data"`
}

type CreateViewRequest struct {
	Data CreateViewData `json:"data"`
}

type CreateViewData struct {
	SearchKeywordURLID int64          `json:"search_keyword_url_id"`
	Name               string         `json:"name"`
	FilterConfig       map[string]any `json:"filter_config"`
}

type GetViewRequest struct {
	ID int64 `json:"id"`
}

type UpdateViewRequest struct {
	ID   int64          `json:"id"`
	Data UpdateViewData `json:"data"`
}

type UpdateViewData struct {
	Name         *string         `json:"name"`
	FilterConfig *map[string]any `json:"filter_config"`
}

type DeleteViewRequest struct {
	ID int64 `json:"id"`
}

type DeleteViewResponse struct {
	Data DeleteViewData `json:"data"`
}

type DeleteViewData struct {
	ID int64 `json:"id"`
}

type ViewPageCountRequest struct {
	ViewID    int64 `json:"view_id"`
	SessionID int64 `json:"crawling_session_id"`
}

type ViewPageCountResponse struct {
	Data ViewPageCountData `json:"data"`
}

type ViewPageCountData struct {
	PageCount int `json:"page_count"`
}
