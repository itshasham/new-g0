package sessions

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	sessionsDto "sitecrawler/newgo/controllers/dto/sessions"
	"sitecrawler/newgo/internal/services/sessions"
	"sitecrawler/newgo/utils/logger"
)

// @Summary Create crawling session
// @Description Creates a new crawling session for a SKU
// @Tags CrawlingSessions
// @Accept json
// @Produce json
// @Param request body sessionsDto.CreateCrawlingSessionRequest true "Crawling session payload"
// @Success 201 {object} sessionsDto.CrawlingSessionResponse
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /api/crawling_sessions [post]
func CreateCrawlingSessionHandler(service sessions.Service) fiber.Handler {
	if service == nil {
		panic("crawling session service required")
	}

	return func(c *fiber.Ctx) error {
		var request sessionsDto.CreateCrawlingSessionRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json payload"})
		}

		if err := validateCreateSessionRequest(request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		resp, err := service.Create(c.Context(), request)
		if err != nil {
			logger.Error(c.UserContext(), "crawling session create failed", logger.Fields{logger.FieldError: err.Error()})
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		if resp.Body == nil {
			return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
		}
		return c.Status(resp.StatusCode).JSON(resp.Body)
	}
}

func validateCreateSessionRequest(req sessionsDto.CreateCrawlingSessionRequest) error {
	if req.Data.SearchKeywordURLID == 0 {
		return errors.New("search_keyword_url_id is required")
	}
	if strings.TrimSpace(req.Data.URL) == "" {
		return errors.New("url is required")
	}
	return nil
}
