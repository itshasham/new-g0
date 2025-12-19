package stats

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"

	statsDto "sitecrawler/newgo/controllers/dto/stats"
	"sitecrawler/newgo/utils/logger"
)

// FetchStats godoc
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
func (ctrl *Controller) Fetch(c *fiber.Ctx) error {
	ctx := c.UserContext()
	fields := logger.Fields{
		logger.FieldMethod: "FetchStats",
	}
	logger.Info(ctx, "stats fetch request received", fields)

	sessionParam := c.Query("crawling_session_id")
	if sessionParam == "" {
		fields[logger.FieldError] = "crawling_session_id is required"
		logger.Error(ctx, "invalid query params", fields)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "crawling_session_id is required"})
	}
	sessionID, err := strconv.ParseInt(sessionParam, 10, 64)
	if err != nil || sessionID == 0 {
		fields[logger.FieldError] = "invalid crawling_session_id"
		logger.Error(ctx, "invalid query params", fields)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid crawling_session_id"})
	}

	var filters []map[string]any
	if rawFilters := c.Query("filters"); rawFilters != "" {
		if err := json.Unmarshal([]byte(rawFilters), &filters); err != nil {
			fields[logger.FieldError] = err.Error()
			logger.Error(ctx, "invalid filters", fields)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid filters"})
		}
	}

	var prefilters []map[string]any
	if rawPrefilters := c.Query("prefilters"); rawPrefilters != "" {
		if err := json.Unmarshal([]byte(rawPrefilters), &prefilters); err != nil {
			fields[logger.FieldError] = err.Error()
			logger.Error(ctx, "invalid prefilters", fields)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid prefilters"})
		}
	}

	var comparisonID *int64
	if rawComparison := c.Query("comparison_crawling_session_id"); rawComparison != "" {
		val, err := strconv.ParseInt(rawComparison, 10, 64)
		if err != nil {
			fields[logger.FieldError] = err.Error()
			logger.Error(ctx, "invalid comparison_crawling_session_id", fields)
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
	fields[logger.FieldRequest] = req
	logger.Info(ctx, "request received", fields)

	resp, err := ctrl.service.Fetch(c.Context(), req)
	if err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "stats fetch failed", fields)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	fields[logger.FieldResponse] = resp
	logger.Info(ctx, "stats fetched successfully", fields)

	if resp.Body == nil {
		return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return c.Status(resp.StatusCode).JSON(resp.Body)
}
