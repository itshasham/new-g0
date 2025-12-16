package models

import "time"

type CrawlingSession struct {
	ID                     int64
	SearchKeywordURLID     int64
	URL                    string
	Status                 string
	Queue                  int
	StartedAt              *time.Time
	EndedAt                *time.Time
	EndReason              string
	Error                  string
	Version                int
	IPs                    []string
	DNSServers             []string
	Aliases                []string
	Location               string
	Sitemap                bool
	Robots                 bool
	SSLValid               bool
	SSLValidUntil          *time.Time
	PagesCount             int
	InternalURLsCount      int
	IgnoredURLsCount       int
	ExternalURLsCount      int
	InternalResourcesCount int
	ExternalResourcesCount int
	Options                map[string]any
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// Health represents the JSON response emitted by /healthz.
type Health struct {
	Status string `json:"status"`
}

type Page struct {
	ID                int64  `json:"id"`
	CrawlingSessionID int64  `json:"crawling_session_id"`
	URL               string `json:"url"`
	ResponseCode      int    `json:"response_code"`
}

type CheckWithPages struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Pages []Page `json:"pages"`
}

type PageImage struct {
	ID     int64  `json:"id"`
	PageID int64  `json:"page_id"`
	URL    string `json:"url"`
}

type AuditCheck struct {
	ID                 int64
	SearchKeywordURLID int64
	Name               string
	Category           string
	FilterConfig       map[string]any
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
