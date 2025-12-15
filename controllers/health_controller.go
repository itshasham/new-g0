package controllers

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/services"
)

type HealthController struct {
	service services.HealthService
	logger  *slog.Logger
}

func NewHealthController(service services.HealthService, logger *slog.Logger) *HealthController {
	if service == nil {
		panic("health service required")
	}
	if logger == nil {
		logger = slog.Default()
	}

	return &HealthController{
		service: service,
		logger:  logger,
	}
}

// @Summary Health check
// @Description Returns service health status
// @Tags Health
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Failure 503 {object} map[string]string
// @Router /healthz [get]
func (c *HealthController) Health(ctx *fiber.Ctx) error {

	status, err := c.service.Status(ctx.Context())
	if err != nil {
		c.logger.Error("health check failed", "error", err)
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "service unavailable"})
	}

	response := dto.HealthResponse{
		Status: status.Status,
	}
	return ctx.JSON(response)
}
