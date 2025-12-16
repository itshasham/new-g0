package routes

import (
	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/controllers"
)

type Dependencies struct {
	Health                *controllers.HealthController
	CrawlingSessionCreate *controllers.CrawlingSessionCreateController
	CrawlingSessionGet    *controllers.CrawlingSessionGetController
	CrawlingSessionPages  *controllers.CrawlingSessionPagesController
	CrawlingSessionChecks *controllers.CrawlingSessionChecksController
}

func Register(app *fiber.App, deps Dependencies) {
	if app == nil {
		panic("fiber app cannot be nil")
	}
	if deps.Health == nil {
		panic("health controller missing")
	}

	app.Get("/healthz", deps.Health.Health)
	if deps.CrawlingSessionCreate != nil {
		app.Post("/api/crawling_sessions", deps.CrawlingSessionCreate.Create)
	}
	if deps.CrawlingSessionGet != nil {
		app.Get("/api/crawling_sessions/:id", deps.CrawlingSessionGet.Get)
	}
	if deps.CrawlingSessionPages != nil {
		app.Get("/api/crawling_sessions/:id/pages", deps.CrawlingSessionPages.List)
	}
	if deps.CrawlingSessionChecks != nil {
		app.Get("/api/crawling_sessions/:id/checks_with_pages", deps.CrawlingSessionChecks.List)
	}
}
