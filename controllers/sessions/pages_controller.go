package sessions

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"

	sessionsDto "sitecrawler/newgo/controllers/dto/sessions"
	"sitecrawler/newgo/internal/services/sessions"
	"sitecrawler/newgo/utils/logger"
)

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
// @Success 200 {object} sessionsDto.CrawlingSessionPagesResponse
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /api/crawling_sessions/{id}/pages [get]
func ListCrawlingSessionPagesHandler(service sessions.Service) fiber.Handler {
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

		page, _ := strconv.Atoi(c.Query("page"))
		pageLimit, _ := strconv.Atoi(c.Query("page_limit"))

		req := sessionsDto.ListCrawlingSessionPagesRequest{
			SessionID: id,
			Filters:   filters,
			Sort:      c.Query("sort"),
			Direction: c.Query("direction"),
			Page:      page,
			PageLimit: pageLimit,
		}

		resp, err := service.ListPages(c.Context(), req)
		if err != nil {
			logger.Error(c.UserContext(), "crawling session pages list failed", logger.Fields{
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
