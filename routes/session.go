package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/controllers/sessions"
	sessionsvc "sitecrawler/newgo/internal/services/sessions"
)

func RegisterCrawlingSessionRoutes(
	ctx context.Context,
	router fiber.Router,
	sessionService sessionsvc.Service,
) {
	_ = ctx

	if sessionService == nil {
		return
	}

	ctrl := sessions.NewController(sessionService)
	router.Post("/", ctrl.Create)
	router.Get("/:id", ctrl.Get)
	router.Get("/:id/pages", ctrl.ListPages)
	router.Get("/:id/checks_with_pages", ctrl.ListChecks)
}
