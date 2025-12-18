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

	router.Get("/", views.ListViewsHandler(viewService))
	router.Post("/", views.CreateViewHandler(viewService))
	router.Get("/:id", views.GetViewHandler(viewService))
	router.Put("/:id", views.UpdateViewHandler(viewService))
	router.Delete("/:id", views.DeleteViewHandler(viewService))
	router.Get("/:id/page_count", views.ViewPageCountHandler(viewService))
}
