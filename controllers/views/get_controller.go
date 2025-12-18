package views

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	viewsDto "sitecrawler/newgo/dto/views"
	"sitecrawler/newgo/internal/services/views"
)

type GetController struct {
	service views.Service
	logger  *slog.Logger
}

func NewGetController(service views.Service, logger *slog.Logger) *GetController {
	if service == nil {
		panic("view get service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &GetController{service: service, logger: logger}
}

func (c *GetController) Get(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	resp, err := c.service.Get(ctx.Context(), viewsDto.GetViewRequest{ID: id})
	if err != nil {
		c.logger.Error("view get failed", "error", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if resp.Body == nil {
		return ctx.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return ctx.Status(resp.StatusCode).JSON(resp.Body)
}
