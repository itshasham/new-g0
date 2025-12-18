package sessions

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"

	sessionsDto "sitecrawler/newgo/controllers/dto/sessions"
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
func (ctrl *Controller) ListPages(c *fiber.Ctx) error {
	ctx := c.UserContext()
	fields := logger.Fields{
		logger.FieldMethod: "ListCrawlingSessionPages",
	}
	logger.Info(ctx, "list crawling session pages request received", fields)

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
	fields[logger.FieldRequest] = req
	logger.Info(ctx, "request received", fields)

	resp, err := ctrl.service.ListPages(c.Context(), req)
	if err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "crawling session pages list failed", fields)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	fields[logger.FieldResponse] = resp
	logger.Info(ctx, "crawling session pages retrieved successfully", fields)

	if resp.Body == nil {
		return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return c.Status(resp.StatusCode).JSON(resp.Body)
}
