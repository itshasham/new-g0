package stats

import (
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsController struct{}

func NewMetricsController() *MetricsController {
	return &MetricsController{}
}

// Metrics godoc
// @Summary Prometheus metrics
// @Description Exposes Prometheus metrics
// @Tags Metrics
// @Produce plain
// @Success 200 {string} string "Metrics in Prometheus text format"
// @Router /metrics [get]
func (c *MetricsController) Metrics(ctx *fiber.Ctx) error {
	return adaptor.HTTPHandler(promhttp.Handler())(ctx)
}
