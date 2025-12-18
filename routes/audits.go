package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/controllers/audits"
	auditsvc "sitecrawler/newgo/internal/services/audits"
)

func RegisterAuditCheckRoutes(
	ctx context.Context,
	router fiber.Router,
	auditService auditsvc.Service,
) {
	_ = ctx

	if auditService == nil {
		return
	}

	router.Get("/", audits.ListAuditChecksHandler(auditService))
	router.Post("/", audits.CreateAuditCheckHandler(auditService))
	router.Get("/:id", audits.GetAuditCheckHandler(auditService))
	router.Put("/:id", audits.UpdateAuditCheckHandler(auditService))
	router.Delete("/:id", audits.DeleteAuditCheckHandler(auditService))
}
