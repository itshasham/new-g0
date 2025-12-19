package audits

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	auditsDto "sitecrawler/newgo/controllers/dto/audits"
	"sitecrawler/newgo/utils/logger"
)

// UpdateAuditCheck godoc
// @Summary Update audit check
// @Description Updates an audit check by ID
// @Tags AuditChecks
// @Accept json
// @Produce json
// @Param id path int true "Audit check ID"
// @Param request body auditsDto.UpdateAuditCheckRequest true "Audit check update payload"
// @Success 200 {object} auditsDto.AuditCheckResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /api/audit_checks/{id} [put]
func (ctrl *Controller) Update(c *fiber.Ctx) error {
	ctx := c.UserContext()
	fields := logger.Fields{
		logger.FieldMethod: "UpdateAuditCheck",
	}
	logger.Info(ctx, "audit check update request received", fields)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		fields[logger.FieldError] = err.Error()
		logger.Error(ctx, "invalid id", fields)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var req auditsDto.UpdateAuditCheckRequest
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
		logger.Error(ctx, "audit check update failed", fields)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	fields[logger.FieldResponse] = resp
	logger.Info(ctx, "audit check updated successfully", fields)

	if resp.Body == nil {
		return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return c.Status(resp.StatusCode).JSON(resp.Body)
}
