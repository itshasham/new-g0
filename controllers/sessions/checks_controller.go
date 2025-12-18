package sessions

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"

	sessionsDto "sitecrawler/newgo/controllers/dto/sessions"
	"sitecrawler/newgo/internal/services/sessions"
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
func ListCrawlingSessionChecksHandler(service sessions.Service) fiber.Handler {
	if service == nil {
		panic("crawling session service required")
	}

	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		var filters []map[string]any
		if rawFilters := c.Query("filters"); rawFilters != "" {
			if err := json.Unmarshal([]byte(rawFilters), &filters); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid filters"})
			}
		}

		var comparisonID *int64
		if rawComparison := c.Query("comparison_crawling_session_id"); rawComparison != "" {
			if comp, err := strconv.ParseInt(rawComparison, 10, 64); err == nil {
				comparisonID = &comp
			} else {
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

		resp, err := service.ListChecks(c.Context(), req)
		if err != nil {
			logger.Error(c.UserContext(), "checks with pages failed", logger.Fields{
				logger.FieldError: err.Error(),
				"id":              id,
			})
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		if resp.Body == nil {
			return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
		}
		return c.Status(resp.StatusCode).JSON(resp.Body)
	}
}
