package controllers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/services"
)

type AuditCheckUpdateController struct {
	service services.AuditCheckUpdateService
	logger  *slog.Logger
}

func NewAuditCheckUpdateController(service services.AuditCheckUpdateService, logger *slog.Logger) *AuditCheckUpdateController {
	if service == nil {
		panic("audit check update service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &AuditCheckUpdateController{service: service, logger: logger}
}

func (c *AuditCheckUpdateController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var request dto.UpdateAuditCheckRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid json payload"})
	}
	request.ID = id

	resp, err := c.service.Update(ctx.Context(), request)
	if err != nil {
		c.logger.Error("audit check update failed", "error", err, "id", id)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if resp.Body == nil {
		return ctx.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return ctx.Status(resp.StatusCode).JSON(resp.Body)
}
