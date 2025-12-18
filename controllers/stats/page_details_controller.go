package stats

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	statsDto "sitecrawler/newgo/controllers/dto/stats"
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
func (ctrl *Controller) Details(c *fiber.Ctx) error {
	ctx := c.UserContext()
	fields := logger.Fields{
		logger.FieldMethod: "PageDetails",
	}
	logger.Info(ctx, "page details request received", fields)

	pageID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "invalid page id", fields)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	skuParam := c.Query("search_keyword_url_id")
	if skuParam == "" {
		fields[logger.FieldError] = "search_keyword_url_id is required"
		logger.Error(ctx, "invalid query params", fields)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "search_keyword_url_id is required"})
	}
	skuID, err := strconv.ParseInt(skuParam, 10, 64)
	if err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "invalid search_keyword_url_id", fields)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid search_keyword_url_id"})
	}

	limit, _ := strconv.Atoi(c.Query("limit"))

	req := statsDto.PageDetailsRequest{
		PageID:             pageID,
		SearchKeywordURLID: skuID,
		Limit:              limit,
	}
	fields[logger.FieldRequest] = req
	logger.Info(ctx, "request received", fields)

	resp, err := ctrl.service.Details(c.Context(), req)
	if err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "page details failed", fields)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	fields[logger.FieldResponse] = resp
	logger.Info(ctx, "page details fetched successfully", fields)

	if resp.Body == nil {
		return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return c.Status(resp.StatusCode).JSON(resp.Body)
}
