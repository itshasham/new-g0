package audits

import (
	"log/slog"
	"net/http"
	"strconv"

	auditsDto "sitecrawler/newgo/dto/audits"
	"sitecrawler/newgo/internal/services/audits"

	"github.com/gofiber/fiber/v2"
)

type DeleteController struct {
	service audits.Service
	logger  *slog.Logger
}

func NewDeleteController(service audits.Service, logger *slog.Logger) *DeleteController {
	if service == nil {
		panic("audit check delete service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &DeleteController{service: service, logger: logger}
}

func (c *DeleteController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := c.service.Delete(ctx.Context(), auditsDto.DeleteAuditCheckRequest{ID: id})
	if err != nil {
		c.logger.Error("audit check delete failed", "error", err, "id", id)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if resp.Body == nil {
		return ctx.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return ctx.Status(resp.StatusCode).JSON(resp.Body)
}
