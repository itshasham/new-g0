package health

import (
	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/controllers/dto"
)

// @Summary Health check
// @Description Returns service health status
// @Tags Health
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Router /healthz [get]
func HealthHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		response := dto.HealthResponse{
			Status: "ok",
		}
		return c.JSON(response)
	}
}
