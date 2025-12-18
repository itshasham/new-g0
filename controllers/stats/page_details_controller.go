package stats

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	statsDto "sitecrawler/newgo/dto/stats"
	"sitecrawler/newgo/internal/services/stats"
)

type PageDetailsController struct {
	service stats.Service
	logger  *slog.Logger
}

func NewPageDetailsController(service stats.Service, logger *slog.Logger) *PageDetailsController {
	if service == nil {
		panic("page details service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &PageDetailsController{
		service: service,
		logger:  logger,
	}
}

// @Summary Get page details
// @Description Returns images or referrers for a page depending on its status
// @Tags Pages
// @Produce json
// @Param id path int true "Page ID"
// @Param search_keyword_url_id query int true "Search keyword URL ID"
// @Param limit query int false "Result limit"
// @Success 200 {object} statsDto.PageDetailsResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/pages/{id}/page_details [get]
func (c *PageDetailsController) Details(ctx *fiber.Ctx) error {
	pageID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	skuParam := ctx.Query("search_keyword_url_id")
	if skuParam == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "search_keyword_url_id is required"})
	}
	skuID, err := strconv.ParseInt(skuParam, 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid search_keyword_url_id"})
	}

	limit, _ := strconv.Atoi(ctx.Query("limit"))

	req := statsDto.PageDetailsRequest{
		PageID:             pageID,
		SearchKeywordURLID: skuID,
		Limit:              limit,
	}

	resp, err := c.service.Details(ctx.Context(), req)
	if err != nil {
		c.logger.Error("page details failed", "error", err, "page_id", pageID)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if resp.Body == nil {
		return ctx.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return ctx.Status(resp.StatusCode).JSON(resp.Body)
}
