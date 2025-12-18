package audits

import "sitecrawler/newgo/models"

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
