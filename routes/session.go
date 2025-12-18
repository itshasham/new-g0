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

	router.Post("/", sessions.CreateCrawlingSessionHandler(sessionService))
	router.Get("/:id", sessions.GetCrawlingSessionHandler(sessionService))
	router.Get("/:id/pages", sessions.ListCrawlingSessionPagesHandler(sessionService))
	router.Get("/:id/checks_with_pages", sessions.ListCrawlingSessionChecksHandler(sessionService))
}
