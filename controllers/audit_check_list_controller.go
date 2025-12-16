package controllers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/services"
)

type AuditCheckListController struct {
	service services.AuditCheckListService
	logger  *slog.Logger
}

func NewAuditCheckListController(service services.AuditCheckListService, logger *slog.Logger) *AuditCheckListController {
	if service == nil {
		panic("audit check list service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &AuditCheckListController{service: service, logger: logger}
}

func (c *AuditCheckListController) List(ctx *fiber.Ctx) error {
	skuID, err := strconv.ParseInt(ctx.Query("search_keyword_url_id"), 10, 64)
	if err != nil || skuID == 0 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": errors.New("missing search_keyword_url_id").Error()})
	}

	resp, err := c.service.List(ctx.Context(), dto.ListAuditChecksRequest{SearchKeywordURLID: skuID})
	if err != nil {
		c.logger.Error("audit check list failed", "error", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if resp.Body == nil {
		return ctx.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return ctx.Status(resp.StatusCode).JSON(resp.Body)
}
