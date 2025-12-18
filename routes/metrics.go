package routes

import (
	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/controllers/stats"
)

func RegisterMetricsRoutes(router fiber.Router, controller *stats.MetricsController) {
	if controller == nil {
		return
	}
	router.Get("/metrics", controller.Metrics)
}
