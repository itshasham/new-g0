package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/controllers/views"
	viewsvc "sitecrawler/newgo/internal/services/views"
)

func RegisterViewRoutes(
	ctx context.Context,
	router fiber.Router,
	viewService viewsvc.Service,
) {
	_ = ctx

	if viewService == nil {
		return
	}

	ctrl := views.NewController(viewService)
	router.Get("/", ctrl.List)
	router.Post("/", ctrl.Create)
	router.Get("/:id", ctrl.Get)
	router.Put("/:id", ctrl.Update)
	router.Delete("/:id", ctrl.Delete)
	router.Get("/:id/page_count", ctrl.PageCount)
}
