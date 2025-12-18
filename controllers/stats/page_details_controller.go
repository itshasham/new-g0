package stats

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	statsDto "sitecrawler/newgo/controllers/dto/stats"
	"sitecrawler/newgo/internal/services/stats"
	"sitecrawler/newgo/utils/logger"
)

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
func PageDetailsHandler(service stats.Service) fiber.Handler {
	if service == nil {
		panic("stats service required")
	}

	return func(c *fiber.Ctx) error {
		pageID, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		skuParam := c.Query("search_keyword_url_id")
		if skuParam == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "search_keyword_url_id is required"})
		}
		skuID, err := strconv.ParseInt(skuParam, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid search_keyword_url_id"})
		}

		limit, _ := strconv.Atoi(c.Query("limit"))

		req := statsDto.PageDetailsRequest{
			PageID:             pageID,
			SearchKeywordURLID: skuID,
			Limit:              limit,
		}

		resp, err := service.Details(c.Context(), req)
		if err != nil {
			logger.Error(c.UserContext(), "page details failed", logger.Fields{
				logger.FieldError: err.Error(),
				"page_id":         pageID,
			})
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		if resp.Body == nil {
			return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
		}
		return c.Status(resp.StatusCode).JSON(resp.Body)
	}
}
