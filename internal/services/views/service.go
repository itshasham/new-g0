package views

import (
"sitecrawler/newgo/dto"
	"context"

	viewsDto "sitecrawler/newgo/dto/views"
)

// Service defines all view operations.
type Service interface {
	List(ctx context.Context, req viewsDto.ListViewsRequest) (*dto.Response[viewsDto.ViewsResponse], error)
	Create(ctx context.Context, req viewsDto.CreateViewRequest) (*dto.Response[viewsDto.ViewResponse], error)
	Get(ctx context.Context, req viewsDto.GetViewRequest) (*dto.Response[viewsDto.ViewResponse], error)
	Update(ctx context.Context, req viewsDto.UpdateViewRequest) (*dto.Response[viewsDto.ViewResponse], error)
	Delete(ctx context.Context, req viewsDto.DeleteViewRequest) (*dto.Response[viewsDto.DeleteViewResponse], error)
	PageCount(ctx context.Context, req viewsDto.ViewPageCountRequest) (*dto.Response[viewsDto.ViewPageCountResponse], error)
}
