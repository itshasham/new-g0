package views

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	viewsDto "sitecrawler/newgo/dto/views"
	"sitecrawler/newgo/internal/services/views"
)

type ListController struct {
	service views.Service
	logger  *slog.Logger
}

func NewListController(service views.Service, logger *slog.Logger) *ListController {
	if service == nil {
		panic("view list service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &ListController{service: service, logger: logger}
}

func (c *ListController) List(ctx *fiber.Ctx) error {
	skuID, _ := strconv.ParseInt(ctx.Query("search_keyword_url_id"), 10, 64)

	resp, err := c.service.List(ctx.Context(), viewsDto.ListViewsRequest{SearchKeywordURLID: skuID})
	if err != nil {
		c.logger.Error("view list failed", "error", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if resp.Body == nil {
		return ctx.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return ctx.Status(resp.StatusCode).JSON(resp.Body)
}
