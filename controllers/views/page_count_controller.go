package views

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"

	viewsDto "sitecrawler/newgo/controllers/dto/views"
	"sitecrawler/newgo/internal/services/views"
	"sitecrawler/newgo/utils/logger"
)

func ViewPageCountHandler(service views.Service) fiber.Handler {
	if service == nil {
		panic("view service required")
	}

	return func(c *fiber.Ctx) error {
		viewID, err := strconv.ParseInt(c.Query("view_id"), 10, 64)
		if err != nil || viewID == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors.New("missing view_id").Error()})
		}

		sessionID, _ := strconv.ParseInt(c.Query("crawling_session_id"), 10, 64)

		resp, err := service.PageCount(c.Context(), viewsDto.ViewPageCountRequest{
			ViewID:    viewID,
			SessionID: sessionID,
		})
		if err != nil {
			logger.Error(c.UserContext(), "view page count failed", logger.Fields{logger.FieldError: err.Error()})
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		if resp.Body == nil {
			return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
		}
		return c.Status(resp.StatusCode).JSON(resp.Body)
	}
}
