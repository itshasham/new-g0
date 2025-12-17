package controllers

import (
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsController struct{}

func NewMetricsController() *MetricsController {
	return &MetricsController{}
}

func (c *MetricsController) Metrics(ctx *fiber.Ctx) error {
	return adaptor.HTTPHandler(promhttp.Handler())(ctx)
}
