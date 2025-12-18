package health

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/dto"
)

// Controller handles health check requests.
type Controller struct {
	logger *slog.Logger
}

// NewController creates a new health controller.
func NewController(logger *slog.Logger) *Controller {
	if logger == nil {
		logger = slog.Default()
	}
	return &Controller{
		logger: logger,
	}
}

// @Summary Health check
// @Description Returns service health status
// @Tags Health
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Router /healthz [get]
func (c *Controller) Health(ctx *fiber.Ctx) error {
	response := dto.HealthResponse{
		Status: "ok",
	}
	return ctx.JSON(response)
}
