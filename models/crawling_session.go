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
