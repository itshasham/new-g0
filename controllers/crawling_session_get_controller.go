package controllers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/services"
)

type CrawlingSessionGetController struct {
	service services.CrawlingSessionGetService
	logger  *slog.Logger
}

func NewCrawlingSessionGetController(service services.CrawlingSessionGetService, logger *slog.Logger) *CrawlingSessionGetController {
	if service == nil {
		panic("crawling session get service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &CrawlingSessionGetController{
		service: service,
		logger:  logger,
	}
}

// @Summary Get crawling session
// @Description Fetches a crawling session by ID
// @Tags CrawlingSessions
// @Produce json
// @Param id path int true "Crawling session ID"
// @Success 200 {object} dto.CrawlingSessionResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/crawling_sessions/{id} [get]
func (c *CrawlingSessionGetController) Get(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	req := dto.GetCrawlingSessionRequest{ID: id}
	resp, err := c.service.Get(ctx.Context(), req)
	if err != nil {
		c.logger.Error("crawling session fetch failed", "error", err, "id", id)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if resp.Body == nil {
		return ctx.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return ctx.Status(resp.StatusCode).JSON(resp.Body)
}
