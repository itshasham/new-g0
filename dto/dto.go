package dto

import "sitecrawler/newgo/models"

type CreateCrawlingSessionRequest struct {
	Data CreateCrawlingSessionData `json:"data"`
}

type CreateCrawlingSessionData struct {
	SearchKeywordURLID int64          `json:"search_keyword_url_id"`
	URL                string         `json:"url"`
	Options            map[string]any `json:"options"`
	Queue              int            `json:"queue"`
}

type CrawlingSessionResponse struct {
	Data models.CrawlingSession `json:"data"`
}

type GetCrawlingSessionRequest struct {
	ID int64 `json:"id"`
}

type GetCrawlingSessionResponse = CrawlingSessionResponse

type ListCrawlingSessionPagesRequest struct {
	SessionID int64            `json:"session_id"`
	Filters   []map[string]any `json:"filters"`
	Sort      string           `json:"sort"`
	Direction string           `json:"direction"`
	Page      int              `json:"page"`
	PageLimit int              `json:"page_limit"`
}

type CrawlingSessionPagesResponse struct {
	Data CrawlingSessionPagesData `json:"data"`
}

type CrawlingSessionPagesData struct {
	Pages      []models.Page `json:"pages"`
	PagesTotal int           `json:"pages_total"`
}

type ListCrawlingSessionChecksRequest struct {
	SessionID           int64            `json:"session_id"`
	ComparisonSessionID *int64           `json:"comparison_session_id"`
	ViewFilters         []map[string]any `json:"view_filters"`
	PageLimitPerCheck   int              `json:"page_limit_per_check"`
}

type CrawlingSessionChecksResponse struct {
	Data CrawlingSessionChecksData `json:"data"`
}

type CrawlingSessionChecksData struct {
	Checks []models.CheckWithPages `json:"checks"`
}

type PageDetailsRequest struct {
	PageID             int64 `json:"page_id"`
	SearchKeywordURLID int64 `json:"search_keyword_url_id"`
	Limit              int   `json:"limit"`
}

type PageDetailsResponse struct {
	Data PageDetailsData `json:"data"`
}

type PageDetailsData struct {
	PageImages              []models.PageImage `json:"page_images,omitempty"`
	BrokenPages             []models.Page      `json:"broken_pages,omitempty"`
	PagesLinkedToBrokenPage []models.Page      `json:"pages_linked_to_broken_page,omitempty"`
}

type StatsRequest struct {
	CrawlingSessionID    int64            `json:"crawling_session_id"`
	Filters              []map[string]any `json:"filters"`
	Prefilters           []map[string]any `json:"prefilters"`
	ComparisonCrawlingID *int64           `json:"comparison_crawling_session_id"`
}

type StatsResponse struct {
	Data map[string]any `json:"data"`
}

// HealthResponse is the DTO returned by the GET /healthz endpoint.
type HealthResponse struct {
	Status string `json:"status"`
}

type ListAuditChecksRequest struct {
	SearchKeywordURLID int64 `json:"search_keyword_url_id"`
}

type AuditCheckResponse struct {
	Data models.AuditCheck `json:"data"`
}

type AuditChecksResponse struct {
	Data []models.AuditCheck `json:"data"`
}

type CreateAuditCheckRequest struct {
	Data CreateAuditCheckData `json:"data"`
}

type CreateAuditCheckData struct {
	SearchKeywordURLID int64          `json:"search_keyword_url_id"`
	Name               string         `json:"name"`
	Category           string         `json:"category"`
	FilterConfig       map[string]any `json:"filter_config"`
}

type GetAuditCheckRequest struct {
	ID int64 `json:"id"`
}

type UpdateAuditCheckRequest struct {
	ID   int64                `json:"id"`
	Data UpdateAuditCheckData `json:"data"`
}

type UpdateAuditCheckData struct {
	Name         *string         `json:"name"`
	Category     *string         `json:"category"`
	FilterConfig *map[string]any `json:"filter_config"`
}

type DeleteAuditCheckRequest struct {
	ID int64 `json:"id"`
}

type DeleteAuditCheckResponse struct {
	Data DeleteAuditCheckData `json:"data"`
}

type DeleteAuditCheckData struct {
	ID int64 `json:"id"`
}

// View DTOs
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
