package audits

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"

	auditsDto "sitecrawler/newgo/controllers/dto/audits"
	"sitecrawler/newgo/internal/services/audits"
	"sitecrawler/newgo/utils/logger"
)

func ListAuditChecksHandler(service audits.Service) fiber.Handler {
	if service == nil {
		panic("audit check service required")
	}

	return func(c *fiber.Ctx) error {
		skuID, err := strconv.ParseInt(c.Query("search_keyword_url_id"), 10, 64)
		if err != nil || skuID == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors.New("missing search_keyword_url_id").Error()})
		}

		resp, err := service.List(c.Context(), auditsDto.ListAuditChecksRequest{SearchKeywordURLID: skuID})
		if err != nil {
			logger.Error(c.UserContext(), "audit check list failed", logger.Fields{logger.FieldError: err.Error()})
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		if resp.Body == nil {
			return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
		}
		return c.Status(resp.StatusCode).JSON(resp.Body)
	}
}
