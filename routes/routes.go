package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/controllers/stats"
	auditsvc "sitecrawler/newgo/internal/services/audits"
	sessionsvc "sitecrawler/newgo/internal/services/sessions"
	statssvc "sitecrawler/newgo/internal/services/stats"
	viewsvc "sitecrawler/newgo/internal/services/views"
)

func RegisterRoutes(
	ctx context.Context,
	app *fiber.App,
	metricsController *stats.MetricsController,
	auditService auditsvc.Service,
	sessionService sessionsvc.Service,
	viewService viewsvc.Service,
	statsService statssvc.Service,
) {
	if app == nil {
		panic("fiber app cannot be nil")
	}

	RegisterSwaggerRoutes(app.Group("/swagger"))
	RegisterHealthCheckRoutes(app.Group("/"))
	RegisterMetricsRoutes(app.Group("/"), metricsController)

	api := app.Group("/api")
	RegisterCrawlingSessionRoutes(ctx, api.Group("/crawling_sessions"), sessionService)
	RegisterAuditCheckRoutes(ctx, api.Group("/audit_checks"), auditService)
	RegisterViewRoutes(ctx, api.Group("/views"), viewService)
	RegisterStatsRoutes(ctx, api.Group("/stats"), statsService)
	RegisterPageDetailsRoutes(ctx, api.Group("/pages"), statsService)
}
