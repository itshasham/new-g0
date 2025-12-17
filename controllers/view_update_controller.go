package controllers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/services"
)

type ViewUpdateController struct {
	service services.ViewUpdateService
	logger  *slog.Logger
}

func NewViewUpdateController(service services.ViewUpdateService, logger *slog.Logger) *ViewUpdateController {
	if service == nil {
		panic("view update service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &ViewUpdateController{service: service, logger: logger}
}

func (c *ViewUpdateController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var request dto.UpdateViewRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid json payload"})
	}
	request.ID = id

	resp, err := c.service.Update(ctx.Context(), request)
	if err != nil {
		c.logger.Error("view update failed", "error", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if resp.Body == nil {
		return ctx.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return ctx.Status(resp.StatusCode).JSON(resp.Body)
}
