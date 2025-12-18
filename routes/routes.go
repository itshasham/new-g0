package routes

import (
	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/controllers/audits"
	"sitecrawler/newgo/controllers/health"
	"sitecrawler/newgo/controllers/sessions"
	"sitecrawler/newgo/controllers/stats"
	"sitecrawler/newgo/controllers/views"
)

type Dependencies struct {
	Health                *health.Controller
	Metrics               *stats.MetricsController
	CrawlingSessionCreate *sessions.CreateController
	CrawlingSessionGet    *sessions.GetController
	CrawlingSessionPages  *sessions.PagesController
	CrawlingSessionChecks *sessions.ChecksController
	PageDetails           *stats.PageDetailsController
	Stats                 *stats.StatsController
	AuditCheckList        *audits.ListController
	AuditCheckCreate      *audits.CreateController
	AuditCheckGet         *audits.GetController
	AuditCheckUpdate      *audits.UpdateController
	AuditCheckDelete      *audits.DeleteController
	ViewList              *views.ListController
	ViewCreate            *views.CreateController
	ViewGet               *views.GetController
	ViewUpdate            *views.UpdateController
	ViewDelete            *views.DeleteController
	ViewPageCount         *views.PageCountController
}

func Register(app *fiber.App, deps Dependencies) {
	if app == nil {
		panic("fiber app cannot be nil")
	}
	if deps.Health == nil {
		panic("health controller missing")
	}

	app.Get("/healthz", deps.Health.Health)
	if deps.Metrics != nil {
		app.Get("/metrics", deps.Metrics.Metrics)
	}
	if deps.CrawlingSessionCreate != nil {
		app.Post("/api/crawling_sessions", deps.CrawlingSessionCreate.Create)
	}
	if deps.CrawlingSessionGet != nil {
		app.Get("/api/crawling_sessions/:id", deps.CrawlingSessionGet.Get)
	}
	if deps.CrawlingSessionPages != nil {
		app.Get("/api/crawling_sessions/:id/pages", deps.CrawlingSessionPages.List)
	}
	if deps.CrawlingSessionChecks != nil {
		app.Get("/api/crawling_sessions/:id/checks_with_pages", deps.CrawlingSessionChecks.List)
	}
	if deps.PageDetails != nil {
		app.Get("/api/pages/:id/page_details", deps.PageDetails.Details)
	}
	if deps.Stats != nil {
		app.Get("/api/stats", deps.Stats.Fetch)
	}

	if deps.AuditCheckList != nil {
		app.Get("/api/audit_checks", deps.AuditCheckList.List)
	}
	if deps.AuditCheckCreate != nil {
		app.Post("/api/audit_checks", deps.AuditCheckCreate.Create)
	}
	if deps.AuditCheckGet != nil {
		app.Get("/api/audit_checks/:id", deps.AuditCheckGet.Get)
	}
	if deps.AuditCheckUpdate != nil {
		app.Put("/api/audit_checks/:id", deps.AuditCheckUpdate.Update)
	}
	if deps.AuditCheckDelete != nil {
		app.Delete("/api/audit_checks/:id", deps.AuditCheckDelete.Delete)
	}

	// View routes
	if deps.ViewList != nil {
		app.Get("/api/views", deps.ViewList.List)
	}
	if deps.ViewCreate != nil {
		app.Post("/api/views", deps.ViewCreate.Create)
	}
	if deps.ViewGet != nil {
		app.Get("/api/views/:id", deps.ViewGet.Get)
	}
	if deps.ViewUpdate != nil {
		app.Put("/api/views/:id", deps.ViewUpdate.Update)
	}
	if deps.ViewDelete != nil {
		app.Delete("/api/views/:id", deps.ViewDelete.Delete)
	}
	if deps.ViewPageCount != nil {
		app.Get("/api/views/:id/page_count", deps.ViewPageCount.PageCount)
	}
}
