package routes

import (
	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/controllers"
)

type Dependencies struct {
	Health                *controllers.HealthController
	CrawlingSessionCreate *controllers.CrawlingSessionCreateController
	CrawlingSessionGet    *controllers.CrawlingSessionGetController
	CrawlingSessionPages  *controllers.CrawlingSessionPagesController
	CrawlingSessionChecks *controllers.CrawlingSessionChecksController
	AuditCheckList        *controllers.AuditCheckListController
	AuditCheckCreate      *controllers.AuditCheckCreateController
	AuditCheckGet         *controllers.AuditCheckGetController
	AuditCheckUpdate      *controllers.AuditCheckUpdateController
	AuditCheckDelete      *controllers.AuditCheckDeleteController
}

func Register(app *fiber.App, deps Dependencies) {
	if app == nil {
		panic("fiber app cannot be nil")
	}
	if deps.Health == nil {
		panic("health controller missing")
	}

	app.Get("/healthz", deps.Health.Health)
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
}
