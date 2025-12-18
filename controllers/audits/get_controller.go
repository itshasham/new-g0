package audits

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	auditsDto "sitecrawler/newgo/controllers/dto/audits"
	"sitecrawler/newgo/internal/services/audits"
	"sitecrawler/newgo/utils/logger"
)

func GetAuditCheckHandler(service audits.Service) fiber.Handler {
	if service == nil {
		panic("audit check service required")
	}

	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		resp, err := service.Get(c.Context(), auditsDto.GetAuditCheckRequest{ID: id})
		if err != nil {
			logger.Error(c.UserContext(), "audit check get failed", logger.Fields{
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
