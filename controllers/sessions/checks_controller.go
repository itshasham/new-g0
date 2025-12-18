package sessions

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"

	sessionsDto "sitecrawler/newgo/controllers/dto/sessions"
	"sitecrawler/newgo/utils/logger"
)

// @Summary List checks with pages
// @Description Lists audit checks with their associated pages for a session
// @Tags CrawlingSessions
// @Produce json
// @Param id path int true "Crawling session ID"
// @Param comparison_crawling_session_id query int false "Comparison session ID"
// @Param filters query string false "View filters JSON"
// @Param page_limit_per_check query int false "Page limit per check"
// @Success 200 {object} sessionsDto.CrawlingSessionChecksResponse
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /api/crawling_sessions/{id}/checks_with_pages [get]
func (ctrl *Controller) ListChecks(c *fiber.Ctx) error {
	ctx := c.UserContext()
	fields := logger.Fields{
		logger.FieldMethod: "ListCrawlingSessionChecks",
	}
	logger.Info(ctx, "list crawling session checks request received", fields)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "invalid session id", fields)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var filters []map[string]any
	if rawFilters := c.Query("filters"); rawFilters != "" {
		if err := json.Unmarshal([]byte(rawFilters), &filters); err != nil {
			fields[logger.FieldError] = err.Error()
			logger.Error(ctx, "invalid filters", fields)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid filters"})
		}
	}

	var comparisonID *int64
	if rawComparison := c.Query("comparison_crawling_session_id"); rawComparison != "" {
		if comp, err := strconv.ParseInt(rawComparison, 10, 64); err == nil {
			comparisonID = &comp
		} else {
			fields[logger.FieldError] = err.Error()
			logger.Error(ctx, "invalid comparison_crawling_session_id", fields)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid comparison_crawling_session_id"})
		}
	}

	pageLimit, _ := strconv.Atoi(c.Query("page_limit_per_check"))

	req := sessionsDto.ListCrawlingSessionChecksRequest{
		SessionID:           id,
		ComparisonSessionID: comparisonID,
		ViewFilters:         filters,
		PageLimitPerCheck:   pageLimit,
	}
	fields[logger.FieldRequest] = req
	logger.Info(ctx, "request received", fields)

	resp, err := ctrl.service.ListChecks(c.Context(), req)
	if err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "checks with pages failed", fields)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	fields[logger.FieldResponse] = resp
	logger.Info(ctx, "checks retrieved successfully", fields)

	if resp.Body == nil {
		return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return c.Status(resp.StatusCode).JSON(resp.Body)
}
