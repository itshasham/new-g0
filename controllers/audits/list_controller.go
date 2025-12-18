package audits

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"

	auditsDto "sitecrawler/newgo/controllers/dto/audits"
	"sitecrawler/newgo/utils/logger"
)

func (ctrl *Controller) List(c *fiber.Ctx) error {
	ctx := c.UserContext()
	fields := logger.Fields{
		logger.FieldMethod: "ListAuditChecks",
	}
	logger.Info(ctx, "audit checks list request received", fields)

	skuID, err := strconv.ParseInt(c.Query("search_keyword_url_id"), 10, 64)
	if err != nil || skuID == 0 {
		fields[logger.FieldError] = "missing search_keyword_url_id"
		logger.Error(ctx, "invalid query params", fields)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors.New("missing search_keyword_url_id").Error()})
	}

	req := auditsDto.ListAuditChecksRequest{SearchKeywordURLID: skuID}
	fields[logger.FieldRequest] = req
	logger.Info(ctx, "request received", fields)

	resp, err := ctrl.service.List(c.Context(), req)
	if err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "audit check list failed", fields)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	fields[logger.FieldResponse] = resp
	logger.Info(ctx, "audit checks retrieved successfully", fields)

	if resp.Body == nil {
		return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return c.Status(resp.StatusCode).JSON(resp.Body)
}
