package views

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	viewsDto "sitecrawler/newgo/controllers/dto/views"
	"sitecrawler/newgo/utils/logger"
)

// ListViews godoc
// @Summary List views
// @Description Lists views for a given `search_keyword_url_id`
// @Tags Views
// @Produce json
// @Param search_keyword_url_id query int false "Search keyword URL ID"
// @Success 200 {object} viewsDto.ViewsResponse
// @Failure 422 {object} map[string]string
// @Router /api/views [get]
func (ctrl *Controller) List(c *fiber.Ctx) error {
	ctx := c.UserContext()
	fields := logger.Fields{
		logger.FieldMethod: "ListViews",
	}
	logger.Info(ctx, "view list request received", fields)

	skuID, _ := strconv.ParseInt(c.Query("search_keyword_url_id"), 10, 64)

	req := viewsDto.ListViewsRequest{SearchKeywordURLID: skuID}
	fields[logger.FieldRequest] = req
	logger.Info(ctx, "request received", fields)

	resp, err := ctrl.service.List(c.Context(), req)
	if err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "view list failed", fields)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	fields[logger.FieldResponse] = resp
	logger.Info(ctx, "views listed successfully", fields)

	if resp.Body == nil {
		return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return c.Status(resp.StatusCode).JSON(resp.Body)
}
