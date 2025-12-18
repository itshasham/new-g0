package views

import (
	"github.com/gofiber/fiber/v2"

	viewsDto "sitecrawler/newgo/controllers/dto/views"
	"sitecrawler/newgo/internal/services/views"
	"sitecrawler/newgo/utils/logger"
)

func CreateViewHandler(service views.Service) fiber.Handler {
	if service == nil {
		panic("view service required")
	}

	return func(c *fiber.Ctx) error {
		var request viewsDto.CreateViewRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		resp, err := service.Create(c.Context(), request)
		if err != nil {
			logger.Error(c.UserContext(), "view create failed", logger.Fields{logger.FieldError: err.Error()})
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		if resp.Body == nil {
			return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
		}
		return c.Status(resp.StatusCode).JSON(resp.Body)
	}
}
