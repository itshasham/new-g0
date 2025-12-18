package audits

import (
	"strconv"

	auditsDto "sitecrawler/newgo/controllers/dto/audits"

	"github.com/gofiber/fiber/v2"
	"sitecrawler/newgo/utils/logger"
)

func (ctrl *Controller) Delete(c *fiber.Ctx) error {
	ctx := c.UserContext()
	fields := logger.Fields{
		logger.FieldMethod: "DeleteAuditCheck",
	}
	logger.Info(ctx, "audit check delete request received", fields)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "invalid id", fields)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	req := auditsDto.DeleteAuditCheckRequest{ID: id}
	fields[logger.FieldRequest] = req
	logger.Info(ctx, "request received", fields)

	resp, err := ctrl.service.Delete(c.Context(), req)
	if err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "audit check delete failed", fields)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	fields[logger.FieldResponse] = resp
	logger.Info(ctx, "audit check deleted successfully", fields)

	if resp.Body == nil {
		return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return c.Status(resp.StatusCode).JSON(resp.Body)
}
