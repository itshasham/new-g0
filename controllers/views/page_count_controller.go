package views

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"

	viewsDto "sitecrawler/newgo/controllers/dto/views"
	"sitecrawler/newgo/utils/logger"
)

func (ctrl *Controller) PageCount(c *fiber.Ctx) error {
	ctx := c.UserContext()
	fields := logger.Fields{
		logger.FieldMethod: "ViewPageCount",
	}
	logger.Info(ctx, "view page count request received", fields)

	viewID, err := strconv.ParseInt(c.Query("view_id"), 10, 64)
	if err != nil || viewID == 0 {
		fields[logger.FieldError] = "missing view_id"
		logger.Error(ctx, "invalid query params", fields)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors.New("missing view_id").Error()})
	}

	sessionID, _ := strconv.ParseInt(c.Query("crawling_session_id"), 10, 64)

	req := viewsDto.ViewPageCountRequest{
		ViewID:    viewID,
		SessionID: sessionID,
	}
	fields[logger.FieldRequest] = req
	logger.Info(ctx, "request received", fields)

	resp, err := ctrl.service.PageCount(c.Context(), req)
	if err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "view page count failed", fields)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	fields[logger.FieldResponse] = resp
	logger.Info(ctx, "view page count fetched successfully", fields)

	if resp.Body == nil {
		return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return c.Status(resp.StatusCode).JSON(resp.Body)
}
