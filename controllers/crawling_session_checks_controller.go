package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/services"
)

type CrawlingSessionChecksController struct {
	service services.CrawlingSessionChecksService
	logger  *slog.Logger
}

func NewCrawlingSessionChecksController(service services.CrawlingSessionChecksService, logger *slog.Logger) *CrawlingSessionChecksController {
	if service == nil {
		panic("crawling session check service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &CrawlingSessionChecksController{
		service: service,
		logger:  logger,
	}
}

// @Summary List checks with pages
// @Description Lists audit checks with their associated pages for a session
// @Tags CrawlingSessions
// @Produce json
// @Param id path int true "Crawling session ID"
// @Param comparison_crawling_session_id query int false "Comparison session ID"
// @Param filters query string false "View filters JSON"
// @Param page_limit_per_check query int false "Page limit per check"
// @Success 200 {object} dto.CrawlingSessionChecksResponse
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /api/crawling_sessions/{id}/checks_with_pages [get]
func (c *CrawlingSessionChecksController) List(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var filters []map[string]any
	if rawFilters := ctx.Query("filters"); rawFilters != "" {
		if err := json.Unmarshal([]byte(rawFilters), &filters); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid filters"})
		}
	}

	var comparisonID *int64
	if rawComparison := ctx.Query("comparison_crawling_session_id"); rawComparison != "" {
		if comp, err := strconv.ParseInt(rawComparison, 10, 64); err == nil {
			comparisonID = &comp
		} else {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid comparison_crawling_session_id"})
		}
	}

	pageLimit, _ := strconv.Atoi(ctx.Query("page_limit_per_check"))

	req := dto.ListCrawlingSessionChecksRequest{
		SessionID:           id,
		ComparisonSessionID: comparisonID,
		ViewFilters:         filters,
		PageLimitPerCheck:   pageLimit,
	}

	resp, err := c.service.List(ctx.Context(), req)
	if err != nil {
		c.logger.Error("checks with pages failed", "error", err, "id", id)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if resp.Body == nil {
		return ctx.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return ctx.Status(resp.StatusCode).JSON(resp.Body)
}
