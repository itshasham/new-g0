package controllers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/services"
)

type AuditCheckDeleteController struct {
	service services.AuditCheckDeleteService
	logger  *slog.Logger
}

func NewAuditCheckDeleteController(service services.AuditCheckDeleteService, logger *slog.Logger) *AuditCheckDeleteController {
	if service == nil {
		panic("audit check delete service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &AuditCheckDeleteController{service: service, logger: logger}
}

func (c *AuditCheckDeleteController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := c.service.Delete(ctx.Context(), dto.DeleteAuditCheckRequest{ID: id})
	if err != nil {
		c.logger.Error("audit check delete failed", "error", err, "id", id)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if resp.Body == nil {
		return ctx.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return ctx.Status(resp.StatusCode).JSON(resp.Body)
}
