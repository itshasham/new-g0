package controllers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/services"
)

type ViewListController struct {
	service services.ViewListService
	logger  *slog.Logger
}

func NewViewListController(service services.ViewListService, logger *slog.Logger) *ViewListController {
	if service == nil {
		panic("view list service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &ViewListController{service: service, logger: logger}
}

func (c *ViewListController) List(ctx *fiber.Ctx) error {
	skuID, _ := strconv.ParseInt(ctx.Query("search_keyword_url_id"), 10, 64)

	resp, err := c.service.List(ctx.Context(), dto.ListViewsRequest{SearchKeywordURLID: skuID})
	if err != nil {
		c.logger.Error("view list failed", "error", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if resp.Body == nil {
		return ctx.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return ctx.Status(resp.StatusCode).JSON(resp.Body)
}
