package stats

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"

	statsDto "sitecrawler/newgo/controllers/dto/stats"
	"sitecrawler/newgo/internal/services/stats"
	"sitecrawler/newgo/utils/logger"
)

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
func FetchStatsHandler(service stats.Service) fiber.Handler {
	if service == nil {
		panic("stats service required")
	}

	return func(c *fiber.Ctx) error {
		sessionParam := c.Query("crawling_session_id")
		if sessionParam == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "crawling_session_id is required"})
		}
		sessionID, err := strconv.ParseInt(sessionParam, 10, 64)
		if err != nil || sessionID == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid crawling_session_id"})
		}

		var filters []map[string]any
		if rawFilters := c.Query("filters"); rawFilters != "" {
			if err := json.Unmarshal([]byte(rawFilters), &filters); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid filters"})
			}
		}

		var prefilters []map[string]any
		if rawPrefilters := c.Query("prefilters"); rawPrefilters != "" {
			if err := json.Unmarshal([]byte(rawPrefilters), &prefilters); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid prefilters"})
			}
		}

		var comparisonID *int64
		if rawComparison := c.Query("comparison_crawling_session_id"); rawComparison != "" {
			val, err := strconv.ParseInt(rawComparison, 10, 64)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid comparison_crawling_session_id"})
			}
			comparisonID = &val
		}

		req := statsDto.StatsRequest{
			CrawlingSessionID:    sessionID,
			Filters:              filters,
			Prefilters:           prefilters,
			ComparisonCrawlingID: comparisonID,
		}

		resp, err := service.Fetch(c.Context(), req)
		if err != nil {
			logger.Error(c.UserContext(), "stats fetch failed", logger.Fields{
				logger.FieldError: err.Error(),
				"session_id":      sessionID,
			})
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		if resp.Body == nil {
			return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
		}
		return c.Status(resp.StatusCode).JSON(resp.Body)
	}
}
