package views

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	viewsDto "sitecrawler/newgo/controllers/dto/views"
	"sitecrawler/newgo/internal/services/views"
	"sitecrawler/newgo/utils/logger"
)

func DeleteViewHandler(service views.Service) fiber.Handler {
	if service == nil {
		panic("view service required")
	}

	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		resp, err := service.Delete(c.Context(), viewsDto.DeleteViewRequest{ID: id})
		if err != nil {
			logger.Error(c.UserContext(), "view delete failed", logger.Fields{logger.FieldError: err.Error()})
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		if resp.Body == nil {
			return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
		}
		return c.Status(resp.StatusCode).JSON(resp.Body)
	}
}
