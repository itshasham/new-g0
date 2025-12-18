package routes

import (
	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/controllers/health"
)

func RegisterHealthCheckRoutes(router fiber.Router) {
	router.Get("/healthz", health.HealthHandler())
}
