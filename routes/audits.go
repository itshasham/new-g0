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

	ctrl := audits.NewController(auditService)
	router.Get("/", ctrl.List)
	router.Post("/", ctrl.Create)
	router.Get("/:id", ctrl.Get)
	router.Put("/:id", ctrl.Update)
	router.Delete("/:id", ctrl.Delete)
}
