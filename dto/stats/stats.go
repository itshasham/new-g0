package stats

import "sitecrawler/newgo/models"

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
