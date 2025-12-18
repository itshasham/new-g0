package sessions

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	sessionsDto "sitecrawler/newgo/controllers/dto/sessions"
	"sitecrawler/newgo/utils/logger"
)

// @Summary Get crawling session
// @Description Fetches a crawling session by ID
// @Tags CrawlingSessions
// @Produce json
// @Param id path int true "Crawling session ID"
// @Success 200 {object} sessionsDto.CrawlingSessionResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/crawling_sessions/{id} [get]
func (ctrl *Controller) Get(c *fiber.Ctx) error {
	ctx := c.UserContext()
	fields := logger.Fields{
		logger.FieldMethod: "GetCrawlingSession",
	}
	logger.Info(ctx, "get crawling session request received", fields)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "invalid id", fields)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	req := sessionsDto.GetCrawlingSessionRequest{ID: id}
	fields[logger.FieldRequest] = req
	logger.Info(ctx, "request received", fields)

	resp, err := ctrl.service.Get(c.Context(), req)
	if err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "crawling session fetch failed", fields)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
	fields[logger.FieldResponse] = resp
	logger.Info(ctx, "crawling session fetched successfully", fields)

	if resp.Body == nil {
		return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return c.Status(resp.StatusCode).JSON(resp.Body)
}
