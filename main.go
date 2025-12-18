package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/controllers/audits"
	"sitecrawler/newgo/controllers/health"
	"sitecrawler/newgo/controllers/sessions"
	"sitecrawler/newgo/controllers/stats"
	"sitecrawler/newgo/controllers/views"
	"sitecrawler/newgo/internal/repository"
	auditsvc "sitecrawler/newgo/internal/services/audits"
	sessionsvc "sitecrawler/newgo/internal/services/sessions"
	statssvc "sitecrawler/newgo/internal/services/stats"
	viewsvc "sitecrawler/newgo/internal/services/views"
	"sitecrawler/newgo/routes"
)

func main() {
	app := fiber.New()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Initialize repositories
	crawlingSessionRepo := repository.NewInMemoryCrawlingSessionRepository()
	pageRepo := repository.NewNoopCrawlingSessionPageRepository()
	checkRepo := repository.NewNoopCrawlingSessionCheckRepository()
	auditRepo := repository.NewInMemoryAuditCheckRepository()
	viewRepo := repository.NewInMemoryViewRepository()
	statsRepo := repository.NewNoopStatsRepository()
	pageDetailsRepo := repository.NewNoopPageDetailsRepository()

	// Health controller
	healthCtrl := health.NewController(logger)

	// Crawling session service and controllers
	sessionSvc := sessionsvc.NewService(crawlingSessionRepo, pageRepo, checkRepo)
	crawlingCreateCtrl := sessions.NewCreateController(sessionSvc, logger)
	crawlingGetCtrl := sessions.NewGetController(sessionSvc, logger)
	crawlingPagesCtrl := sessions.NewPagesController(sessionSvc, logger)
	crawlingChecksCtrl := sessions.NewChecksController(sessionSvc, logger)

	// Audit check service and controllers
	auditSvc := auditsvc.NewService(auditRepo)
	auditListCtrl := audits.NewListController(auditSvc, logger)
	auditCreateCtrl := audits.NewCreateController(auditSvc, logger)
	auditGetCtrl := audits.NewGetController(auditSvc, logger)
	auditUpdateCtrl := audits.NewUpdateController(auditSvc, logger)
	auditDeleteCtrl := audits.NewDeleteController(auditSvc, logger)

	// View service and controllers
	viewSvc := viewsvc.NewService(viewRepo, pageRepo)
	viewListCtrl := views.NewListController(viewSvc, logger)
	viewCreateCtrl := views.NewCreateController(viewSvc, logger)
	viewGetCtrl := views.NewGetController(viewSvc, logger)
	viewUpdateCtrl := views.NewUpdateController(viewSvc, logger)
	viewDeleteCtrl := views.NewDeleteController(viewSvc, logger)
	viewPageCountCtrl := views.NewPageCountController(viewSvc, logger)

	// Stats service and controllers
	statsSvc := statssvc.NewService(statsRepo, pageDetailsRepo)
	metricsCtrl := stats.NewMetricsController()
	statsCtrl := stats.NewStatsController(statsSvc, logger)
	pageDetailsCtrl := stats.NewPageDetailsController(statsSvc, logger)

	routes.Register(app, routes.Dependencies{
		Health:                healthCtrl,
		Metrics:               metricsCtrl,
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
