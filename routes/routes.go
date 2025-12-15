package routes

import (
	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/controllers"
)

type Dependencies struct {
	Health           *controllers.HealthController
	CrawlingSessions *controllers.CrawlingSessionController
}

func Register(app *fiber.App, deps Dependencies) {
	if app == nil {
		panic("fiber app cannot be nil")
	}
	if deps.Health == nil {
		panic("health controller missing")
	}

	app.Get("/healthz", deps.Health.Health)
	if deps.CrawlingSessions != nil {
		app.Post("/api/crawling_sessions", deps.CrawlingSessions.Create)
	}
}
