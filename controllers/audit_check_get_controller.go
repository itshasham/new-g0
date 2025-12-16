package controllers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/services"
)

type AuditCheckGetController struct {
	service services.AuditCheckGetService
	logger  *slog.Logger
}

func NewAuditCheckGetController(service services.AuditCheckGetService, logger *slog.Logger) *AuditCheckGetController {
	if service == nil {
		panic("audit check get service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &AuditCheckGetController{service: service, logger: logger}
}

func (c *AuditCheckGetController) Get(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := c.service.Get(ctx.Context(), dto.GetAuditCheckRequest{ID: id})
	if err != nil {
		c.logger.Error("audit check get failed", "error", err, "id", id)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if resp.Body == nil {
		return ctx.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return ctx.Status(resp.StatusCode).JSON(resp.Body)
}
