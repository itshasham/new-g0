package views

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	viewsDto "sitecrawler/newgo/controllers/dto/views"
	"sitecrawler/newgo/internal/services/views"
	"sitecrawler/newgo/utils/logger"
)

func ListViewsHandler(service views.Service) fiber.Handler {
	if service == nil {
		panic("view service required")
	}

	return func(c *fiber.Ctx) error {
		skuID, _ := strconv.ParseInt(c.Query("search_keyword_url_id"), 10, 64)

		resp, err := service.List(c.Context(), viewsDto.ListViewsRequest{SearchKeywordURLID: skuID})
		if err != nil {
			logger.Error(c.UserContext(), "view list failed", logger.Fields{logger.FieldError: err.Error()})
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		if resp.Body == nil {
			return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
		}
		return c.Status(resp.StatusCode).JSON(resp.Body)
	}
}
