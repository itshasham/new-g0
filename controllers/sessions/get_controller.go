package sessions

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	sessionsDto "sitecrawler/newgo/controllers/dto/sessions"
	"sitecrawler/newgo/internal/services/sessions"
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
func GetCrawlingSessionHandler(service sessions.Service) fiber.Handler {
	if service == nil {
		panic("crawling session service required")
	}

	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		req := sessionsDto.GetCrawlingSessionRequest{ID: id}
		resp, err := service.Get(c.Context(), req)
		if err != nil {
			logger.Error(c.UserContext(), "crawling session fetch failed", logger.Fields{
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
