package sessions

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"

	sessionsDto "sitecrawler/newgo/dto/sessions"
	"sitecrawler/newgo/internal/services/sessions"
)

type CreateController struct {
	service sessions.Service
	logger  *slog.Logger
}

func NewCreateController(service sessions.Service, logger *slog.Logger) *CreateController {
	if service == nil {
		panic("crawling session create service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &CreateController{
		service: service,
		logger:  logger,
	}
}

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
func (c *CreateController) Create(ctx *fiber.Ctx) error {
	var request sessionsDto.CreateCrawlingSessionRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid json payload"})
	}

	if err := validateCreateSessionRequest(request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := c.service.Create(ctx.Context(), request)
	if err != nil {
		c.logger.Error("crawling session create failed", "error", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if resp.Body == nil {
		return ctx.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return ctx.Status(resp.StatusCode).JSON(resp.Body)
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
