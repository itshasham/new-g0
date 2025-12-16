package controllers

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/services"
)

type AuditCheckCreateController struct {
	service services.AuditCheckCreateService
	logger  *slog.Logger
}

func NewAuditCheckCreateController(service services.AuditCheckCreateService, logger *slog.Logger) *AuditCheckCreateController {
	if service == nil {
		panic("audit check create service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &AuditCheckCreateController{service: service, logger: logger}
}

func (c *AuditCheckCreateController) Create(ctx *fiber.Ctx) error {
	var request dto.CreateAuditCheckRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid json payload"})
	}

	if err := validateCreateAuditCheckRequest(request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := c.service.Create(ctx.Context(), request)
	if err != nil {
		c.logger.Error("audit check create failed", "error", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if resp.Body == nil {
		return ctx.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return ctx.Status(resp.StatusCode).JSON(resp.Body)
}

func validateCreateAuditCheckRequest(req dto.CreateAuditCheckRequest) error {
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
