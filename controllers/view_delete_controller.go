package controllers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/services"
)

type ViewDeleteController struct {
	service services.ViewDeleteService
	logger  *slog.Logger
}

func NewViewDeleteController(service services.ViewDeleteService, logger *slog.Logger) *ViewDeleteController {
	if service == nil {
		panic("view delete service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &ViewDeleteController{service: service, logger: logger}
}

func (c *ViewDeleteController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	resp, err := c.service.Delete(ctx.Context(), dto.DeleteViewRequest{ID: id})
	if err != nil {
		c.logger.Error("view delete failed", "error", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if resp.Body == nil {
		return ctx.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return ctx.Status(resp.StatusCode).JSON(resp.Body)
}
