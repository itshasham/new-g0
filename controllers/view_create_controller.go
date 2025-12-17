package controllers

import (
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/services"
)

type ViewCreateController struct {
	service services.ViewCreateService
	logger  *slog.Logger
}

func NewViewCreateController(service services.ViewCreateService, logger *slog.Logger) *ViewCreateController {
	if service == nil {
		panic("view create service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &ViewCreateController{service: service, logger: logger}
}

func (c *ViewCreateController) Create(ctx *fiber.Ctx) error {
	var request dto.CreateViewRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := c.service.Create(ctx.Context(), request)
	if err != nil {
		c.logger.Error("view create failed", "error", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if resp.Body == nil {
		return ctx.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return ctx.Status(resp.StatusCode).JSON(resp.Body)
}
