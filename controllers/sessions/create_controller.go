package sessions

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	sessionsDto "sitecrawler/newgo/controllers/dto/sessions"
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
func (ctrl *Controller) Create(c *fiber.Ctx) error {
	ctx := c.UserContext()
	fields := logger.Fields{
		logger.FieldMethod: "CreateCrawlingSession",
	}
	logger.Info(ctx, "create crawling session request received", fields)

	var req sessionsDto.CreateCrawlingSessionRequest
	if err := c.BodyParser(&req); err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "failed to parse request body", fields)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json payload"})
	}
	fields[logger.FieldRequest] = req
	logger.Info(ctx, "request received", fields)

	if err := validateCreateSessionRequest(req); err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "validation failed", fields)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := ctrl.service.Create(c.Context(), req)
	if err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "crawling session create failed", fields)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
	fields[logger.FieldResponse] = resp
	logger.Info(ctx, "crawling session created successfully", fields)

	if resp.Body == nil {
		return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return c.Status(resp.StatusCode).JSON(resp.Body)
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
