package controllers

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/services"
)

type CrawlingSessionController struct {
	service services.CrawlingSessionService
	logger  *slog.Logger
}

func NewCrawlingSessionController(service services.CrawlingSessionService, logger *slog.Logger) *CrawlingSessionController {
	if service == nil {
		panic("crawling session service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &CrawlingSessionController{
		service: service,
		logger:  logger,
	}
}

// @Summary Create crawling session
// @Description Creates a new crawling session for a SKU
// @Tags CrawlingSessions
// @Accept json
// @Produce json
// @Param request body dto.CreateCrawlingSessionRequest true "Crawling session payload"
// @Success 201 {object} dto.CrawlingSessionResponse
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /api/crawling_sessions [post]
func (c *CrawlingSessionController) Create(ctx *fiber.Ctx) error {

	var request dto.CreateCrawlingSessionRequest
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

func validateCreateSessionRequest(req dto.CreateCrawlingSessionRequest) error {
	if req.Data.SearchKeywordURLID == 0 {
		return errors.New("search_keyword_url_id is required")
	}
	if strings.TrimSpace(req.Data.URL) == "" {
		return errors.New("url is required")
	}
	return nil
}
