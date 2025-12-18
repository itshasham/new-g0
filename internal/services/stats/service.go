package stats

import (
"sitecrawler/newgo/dto"
	"context"

	statsDto "sitecrawler/newgo/dto/stats"
)

// Service defines all stats operations.
type Service interface {
	Fetch(ctx context.Context, req statsDto.StatsRequest) (*dto.Response[statsDto.StatsResponse], error)
	Details(ctx context.Context, req statsDto.PageDetailsRequest) (*dto.Response[statsDto.PageDetailsResponse], error)
}
