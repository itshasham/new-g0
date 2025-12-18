package audits

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	auditsDto "sitecrawler/newgo/controllers/dto/audits"
	"sitecrawler/newgo/internal/services/audits"
	"sitecrawler/newgo/utils/logger"
)

func UpdateAuditCheckHandler(service audits.Service) fiber.Handler {
	if service == nil {
		panic("audit check service required")
	}

	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		var request auditsDto.UpdateAuditCheckRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json payload"})
		}
		request.ID = id

		resp, err := service.Update(c.Context(), request)
		if err != nil {
			logger.Error(c.UserContext(), "audit check update failed", logger.Fields{
				logger.FieldError: err.Error(),
				"id":              id,
			})
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		if resp.Body == nil {
			return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
		}
		return c.Status(resp.StatusCode).JSON(resp.Body)
	}
}
