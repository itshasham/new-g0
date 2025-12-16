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

type CrawlingSessionPagesController struct {
	service services.CrawlingSessionPagesService
	logger  *slog.Logger
}

func NewCrawlingSessionPagesController(service services.CrawlingSessionPagesService, logger *slog.Logger) *CrawlingSessionPagesController {
	if service == nil {
		panic("crawling session page service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &CrawlingSessionPagesController{
		service: service,
		logger:  logger,
	}
}

// @Summary List pages for crawling session
// @Description Lists pages for a crawling session with optional filters
// @Tags CrawlingSessions
// @Produce json
// @Param id path int true "Crawling session ID"
// @Param filters query string false "JSON encoded filters"
// @Param sort query string false "Sort field"
// @Param direction query string false "Sort direction"
// @Param page query int false "Page number"
// @Param page_limit query int false "Page size"
// @Success 200 {object} dto.CrawlingSessionPagesResponse
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /api/crawling_sessions/{id}/pages [get]
func (c *CrawlingSessionPagesController) List(ctx *fiber.Ctx) error {
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

	page, _ := strconv.Atoi(ctx.Query("page"))
	pageLimit, _ := strconv.Atoi(ctx.Query("page_limit"))

	req := dto.ListCrawlingSessionPagesRequest{
		SessionID: id,
		Filters:   filters,
		Sort:      ctx.Query("sort"),
		Direction: ctx.Query("direction"),
		Page:      page,
		PageLimit: pageLimit,
	}

	resp, err := c.service.List(ctx.Context(), req)
	if err != nil {
		c.logger.Error("crawling session pages list failed", "error", err, "id", id)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if resp.Body == nil {
		return ctx.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return ctx.Status(resp.StatusCode).JSON(resp.Body)
}
