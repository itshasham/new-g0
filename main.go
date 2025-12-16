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
	)
	healthCtrl := controllers.NewHealthController(svc.Health(), logger)
	crawlingCreateCtrl := controllers.NewCrawlingSessionCreateController(svc.CrawlingSessionCreator(), logger)
	crawlingGetCtrl := controllers.NewCrawlingSessionGetController(svc.CrawlingSessionGetter(), logger)
	crawlingPagesCtrl := controllers.NewCrawlingSessionPagesController(svc.CrawlingSessionPages(), logger)
	crawlingChecksCtrl := controllers.NewCrawlingSessionChecksController(svc.CrawlingSessionChecks(), logger)
	auditListCtrl := controllers.NewAuditCheckListController(svc.AuditCheckLister(), logger)
	auditCreateCtrl := controllers.NewAuditCheckCreateController(svc.AuditCheckCreator(), logger)
	auditGetCtrl := controllers.NewAuditCheckGetController(svc.AuditCheckGetter(), logger)
	auditUpdateCtrl := controllers.NewAuditCheckUpdateController(svc.AuditCheckUpdater(), logger)
	auditDeleteCtrl := controllers.NewAuditCheckDeleteController(svc.AuditCheckDeleter(), logger)

	routes.Register(app, routes.Dependencies{
		Health:                healthCtrl,
		CrawlingSessionCreate: crawlingCreateCtrl,
		CrawlingSessionGet:    crawlingGetCtrl,
		CrawlingSessionPages:  crawlingPagesCtrl,
		CrawlingSessionChecks: crawlingChecksCtrl,
		AuditCheckList:        auditListCtrl,
		AuditCheckCreate:      auditCreateCtrl,
		AuditCheckGet:         auditGetCtrl,
		AuditCheckUpdate:      auditUpdateCtrl,
		AuditCheckDelete:      auditDeleteCtrl,
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
