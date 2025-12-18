package stats

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	statsDto "sitecrawler/newgo/dto/stats"
	"sitecrawler/newgo/internal/services/stats"
)

type StatsController struct {
	service stats.Service
	logger  *slog.Logger
}

func NewStatsController(service stats.Service, logger *slog.Logger) *StatsController {
	if service == nil {
		panic("stats service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &StatsController{
		service: service,
		logger:  logger,
	}
}

// @Summary Fetch stats
// @Description Retrieves stats for a crawling session with optional filters
// @Tags Stats
// @Produce json
// @Param crawling_session_id query int true "Crawling session ID"
// @Param filters query string false "JSON filters"
// @Param prefilters query string false "JSON prefilters"
// @Param comparison_crawling_session_id query int false "Comparison session ID"
// @Success 200 {object} statsDto.StatsResponse
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /api/stats [get]
func (c *StatsController) Fetch(ctx *fiber.Ctx) error {
	sessionParam := ctx.Query("crawling_session_id")
	if sessionParam == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "crawling_session_id is required"})
	}
	sessionID, err := strconv.ParseInt(sessionParam, 10, 64)
	if err != nil || sessionID == 0 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid crawling_session_id"})
	}

	var filters []map[string]any
	if rawFilters := ctx.Query("filters"); rawFilters != "" {
		if err := json.Unmarshal([]byte(rawFilters), &filters); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid filters"})
		}
	}

	var prefilters []map[string]any
	if rawPrefilters := ctx.Query("prefilters"); rawPrefilters != "" {
		if err := json.Unmarshal([]byte(rawPrefilters), &prefilters); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid prefilters"})
		}
	}

	var comparisonID *int64
	if rawComparison := ctx.Query("comparison_crawling_session_id"); rawComparison != "" {
		val, err := strconv.ParseInt(rawComparison, 10, 64)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid comparison_crawling_session_id"})
		}
		comparisonID = &val
	}

	req := statsDto.StatsRequest{
		CrawlingSessionID:    sessionID,
		Filters:              filters,
		Prefilters:           prefilters,
		ComparisonCrawlingID: comparisonID,
	}

	resp, err := c.service.Fetch(ctx.Context(), req)
	if err != nil {
		c.logger.Error("stats fetch failed", "error", err, "session_id", sessionID)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if resp.Body == nil {
		return ctx.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return ctx.Status(resp.StatusCode).JSON(resp.Body)
}
