package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/controllers"
	"sitecrawler/newgo/internal/repository"
	"sitecrawler/newgo/internal/services"
	"sitecrawler/newgo/routes"
)

func main() {
	app := fiber.New()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	svc := services.New(
		services.WithHealthRepository(repository.NewNoopHealthRepository()),
		services.WithCrawlingSessionRepository(repository.NewInMemoryCrawlingSessionRepository()),
		services.WithAuditCheckRepository(repository.NewInMemoryAuditCheckRepository()),
		services.WithViewRepository(repository.NewInMemoryViewRepository()),
	)
	healthCtrl := controllers.NewHealthController(svc.Health(), logger)
	crawlingCreateCtrl := controllers.NewCrawlingSessionCreateController(svc.CrawlingSessionCreator(), logger)
	crawlingGetCtrl := controllers.NewCrawlingSessionGetController(svc.CrawlingSessionGetter(), logger)
	crawlingPagesCtrl := controllers.NewCrawlingSessionPagesController(svc.CrawlingSessionPages(), logger)
	crawlingChecksCtrl := controllers.NewCrawlingSessionChecksController(svc.CrawlingSessionChecks(), logger)
	pageDetailsCtrl := controllers.NewPageDetailsController(svc.PageDetails(), logger)
	statsCtrl := controllers.NewStatsController(svc.Stats(), logger)
	auditListCtrl := controllers.NewAuditCheckListController(svc.AuditCheckLister(), logger)
	auditCreateCtrl := controllers.NewAuditCheckCreateController(svc.AuditCheckCreator(), logger)
	auditGetCtrl := controllers.NewAuditCheckGetController(svc.AuditCheckGetter(), logger)
	auditUpdateCtrl := controllers.NewAuditCheckUpdateController(svc.AuditCheckUpdater(), logger)
	auditDeleteCtrl := controllers.NewAuditCheckDeleteController(svc.AuditCheckDeleter(), logger)
	viewListCtrl := controllers.NewViewListController(svc.ViewLister(), logger)
	viewCreateCtrl := controllers.NewViewCreateController(svc.ViewCreator(), logger)
	viewGetCtrl := controllers.NewViewGetController(svc.ViewGetter(), logger)
	viewUpdateCtrl := controllers.NewViewUpdateController(svc.ViewUpdater(), logger)
	viewDeleteCtrl := controllers.NewViewDeleteController(svc.ViewDeleter(), logger)
	viewPageCountCtrl := controllers.NewViewPageCountController(svc.ViewPageCounter(), logger)

	routes.Register(app, routes.Dependencies{
		Health:                healthCtrl,
		CrawlingSessionCreate: crawlingCreateCtrl,
		CrawlingSessionGet:    crawlingGetCtrl,
		CrawlingSessionPages:  crawlingPagesCtrl,
		CrawlingSessionChecks: crawlingChecksCtrl,
		PageDetails:           pageDetailsCtrl,
		Stats:                 statsCtrl,
		AuditCheckList:        auditListCtrl,
		AuditCheckCreate:      auditCreateCtrl,
		AuditCheckGet:         auditGetCtrl,
		AuditCheckUpdate:      auditUpdateCtrl,
		AuditCheckDelete:      auditDeleteCtrl,
		ViewList:              viewListCtrl,
		ViewCreate:            viewCreateCtrl,
		ViewGet:               viewGetCtrl,
		ViewUpdate:            viewUpdateCtrl,
		ViewDelete:            viewDeleteCtrl,
		ViewPageCount:         viewPageCountCtrl,
	})

	addr := getenv("ADDR", ":8080")
	startServer(app, addr, logger)
}

func startServer(app *fiber.App, addr string, logger *slog.Logger) {
	go func() {
		logger.Info("fiber server starting", "addr", addr)
		if err := app.Listen(addr); err != nil {
			logger.Error("fiber server stopped", "error", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("shutdown signal received")
	if err := app.Shutdown(); err != nil {
		logger.Error("fiber shutdown error", "error", err)
		return
	}

	log.Println("fiber server stopped gracefully")
}

func getenv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return fallback
}
