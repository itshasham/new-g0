package audits

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	auditsDto "sitecrawler/newgo/controllers/dto/audits"

	"sitecrawler/newgo/internal/services/audits"
	"sitecrawler/newgo/utils/logger"
)

func CreateAuditCheckHandler(service audits.Service) fiber.Handler {
	if service == nil {
		panic("audit check service required")
	}

	return func(c *fiber.Ctx) error {
		var request auditsDto.CreateAuditCheckRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json payload"})
		}

		if err := validateCreateAuditCheckRequest(request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		resp, err := service.Create(c.Context(), request)
		if err != nil {
			logger.Error(c.UserContext(), "audit check create failed", logger.Fields{logger.FieldError: err.Error()})
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		if resp.Body == nil {
			return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
		}
		return c.Status(resp.StatusCode).JSON(resp.Body)
	}
}

func validateCreateAuditCheckRequest(req auditsDto.CreateAuditCheckRequest) error {
	if req.Data.SearchKeywordURLID == 0 {
		return errors.New("search_keyword_url_id is required")
	}
	if strings.TrimSpace(req.Data.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(req.Data.Category) == "" {
		return errors.New("category is required")
	}
	return nil
}
