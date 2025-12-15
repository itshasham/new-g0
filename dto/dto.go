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

// HealthResponse is the DTO returned by the GET /healthz endpoint.
type HealthResponse struct {
	Status string `json:"status"`
}
