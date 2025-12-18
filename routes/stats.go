package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/controllers/stats"
	statssvc "sitecrawler/newgo/internal/services/stats"
)

func RegisterStatsRoutes(ctx context.Context, router fiber.Router, statsService statssvc.Service) {
	_ = ctx

	if statsService == nil {
		return
	}
	ctrl := stats.NewController(statsService)
	router.Get("/", ctrl.Fetch)
}

func RegisterPageDetailsRoutes(ctx context.Context, router fiber.Router, statsService statssvc.Service) {
	_ = ctx

	if statsService == nil {
		return
	}
	ctrl := stats.NewController(statsService)
	router.Get("/:id/page_details", ctrl.Details)
}
