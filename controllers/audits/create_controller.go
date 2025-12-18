package audits

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	auditsDto "sitecrawler/newgo/controllers/dto/audits"

	"sitecrawler/newgo/utils/logger"
)

func (ctrl *Controller) Create(c *fiber.Ctx) error {
	ctx := c.UserContext()
	fields := logger.Fields{
		logger.FieldMethod: "CreateAuditCheck",
	}
	logger.Info(ctx, "audit check create request received", fields)

	var req auditsDto.CreateAuditCheckRequest
	if err := c.BodyParser(&req); err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "failed to parse request body", fields)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json payload"})
	}
	fields[logger.FieldRequest] = req
	logger.Info(ctx, "request received", fields)

	if err := validateCreateAuditCheckRequest(req); err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "validation failed", fields)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := ctrl.service.Create(c.Context(), req)
	if err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "audit check create failed", fields)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	fields[logger.FieldResponse] = resp
	logger.Info(ctx, "audit check created successfully", fields)

	if resp.Body == nil {
		return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return c.Status(resp.StatusCode).JSON(resp.Body)
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
