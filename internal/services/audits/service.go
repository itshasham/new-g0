package audits

import (
	"context"
	"sitecrawler/newgo/dto"

	auditsDto "sitecrawler/newgo/dto/audits"
)

// Service defines all audit check operations.
type Service interface {
	List(ctx context.Context, req auditsDto.ListAuditChecksRequest) (*dto.Response[auditsDto.AuditChecksResponse], error)
	Create(ctx context.Context, req auditsDto.CreateAuditCheckRequest) (*dto.Response[auditsDto.AuditCheckResponse], error)
	Get(ctx context.Context, req auditsDto.GetAuditCheckRequest) (*dto.Response[auditsDto.AuditCheckResponse], error)
	Update(ctx context.Context, req auditsDto.UpdateAuditCheckRequest) (*dto.Response[auditsDto.AuditCheckResponse], error)
	Delete(ctx context.Context, req auditsDto.DeleteAuditCheckRequest) (*dto.Response[auditsDto.DeleteAuditCheckResponse], error)
}
