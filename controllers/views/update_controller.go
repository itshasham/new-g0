package views

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	viewsDto "sitecrawler/newgo/dto/views"
	"sitecrawler/newgo/internal/services/views"
)

type UpdateController struct {
	service views.Service
	logger  *slog.Logger
}

func NewUpdateController(service views.Service, logger *slog.Logger) *UpdateController {
	if service == nil {
		panic("view update service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &UpdateController{service: service, logger: logger}
}

func (c *UpdateController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var request viewsDto.UpdateViewRequest
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
