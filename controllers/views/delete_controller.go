package views

import (
viewsDto "sitecrawler/newgo/dto/views"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/internal/services/views"
)

type DeleteController struct {
	service views.Service
	logger  *slog.Logger
}

func NewDeleteController(service views.Service, logger *slog.Logger) *DeleteController {
	if service == nil {
		panic("view delete service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &DeleteController{service: service, logger: logger}
}

func (c *DeleteController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	resp, err := c.service.Delete(ctx.Context(), viewsDto.DeleteViewRequest{ID: id})
	if err != nil {
		c.logger.Error("view delete failed", "error", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if resp.Body == nil {
		return ctx.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return ctx.Status(resp.StatusCode).JSON(resp.Body)
}
