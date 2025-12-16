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
