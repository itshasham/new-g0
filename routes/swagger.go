package routes

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/gofiber/swagger"

	// Auto-generated swagger docs
	_ "sitecrawler/newgo/docs"
)

func RegisterSwaggerRoutes(router fiber.Router) {
	// Swagger UI handler
	router.Get("/*", fiberSwagger.HandlerDefault)

	// Info endpoint for when docs aren't generated yet
	router.Get("/info", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message":  "To generate swagger docs, run: make swagger (or: go generate ./...)",
			"docs_url": "/swagger/",
		})
	})
}
