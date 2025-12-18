package views

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	viewsDto "sitecrawler/newgo/controllers/dto/views"
	"sitecrawler/newgo/utils/logger"
)

func (ctrl *Controller) Update(c *fiber.Ctx) error {
	ctx := c.UserContext()
	fields := logger.Fields{
		logger.FieldMethod: "UpdateView",
	}
	logger.Info(ctx, "view update request received", fields)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "invalid id", fields)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var req viewsDto.UpdateViewRequest
	if err := c.BodyParser(&req); err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "failed to parse request body", fields)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json payload"})
	}
	req.ID = id
	fields[logger.FieldRequest] = req
	logger.Info(ctx, "request received", fields)

	resp, err := ctrl.service.Update(c.Context(), req)
	if err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "view update failed", fields)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	fields[logger.FieldResponse] = resp
	logger.Info(ctx, "view updated successfully", fields)

	if resp.Body == nil {
		return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return c.Status(resp.StatusCode).JSON(resp.Body)
}
