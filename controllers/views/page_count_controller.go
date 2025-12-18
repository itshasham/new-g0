package views

import (
viewsDto "sitecrawler/newgo/dto/views"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/internal/services/views"
)

type PageCountController struct {
	service views.Service
	logger  *slog.Logger
}

func NewPageCountController(service views.Service, logger *slog.Logger) *PageCountController {
	if service == nil {
		panic("view page count service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &PageCountController{service: service, logger: logger}
}

func (c *PageCountController) PageCount(ctx *fiber.Ctx) error {
	viewID, err := strconv.ParseInt(ctx.Query("view_id"), 10, 64)
	if err != nil || viewID == 0 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": errors.New("missing view_id").Error()})
	}

	sessionID, _ := strconv.ParseInt(ctx.Query("crawling_session_id"), 10, 64)

	resp, err := c.service.PageCount(ctx.Context(), viewsDto.ViewPageCountRequest{
		ViewID:    viewID,
		SessionID: sessionID,
	})
	if err != nil {
		c.logger.Error("view page count failed", "error", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if resp.Body == nil {
		return ctx.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return ctx.Status(resp.StatusCode).JSON(resp.Body)
}
