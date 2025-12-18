// @title SiteCrawler API
// @version 1.0
// @description SiteCrawler API server
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-Access-Token
// @description Access token for authentication
package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"sitecrawler/newgo/config"
	"sitecrawler/newgo/controllers/stats"
	"sitecrawler/newgo/internal/repository"
	"sitecrawler/newgo/internal/repository/postgres"
	auditsvc "sitecrawler/newgo/internal/services/audits"
	sessionsvc "sitecrawler/newgo/internal/services/sessions"
	statssvc "sitecrawler/newgo/internal/services/stats"
	viewsvc "sitecrawler/newgo/internal/services/views"
	"sitecrawler/newgo/routes"
	"sitecrawler/newgo/storage/database"
	"sitecrawler/newgo/utils/logger"
)

func main() {
	ctx := context.Background()

	config.LoadEnvVariables()
	env := config.LoadEnvConfiguration()

	ctx = logger.WithCorrelationID(ctx, "sitecrawler")
	logger.Info(ctx, "Starting application", logger.Fields{})

	db, err := database.InitDatabase(ctx, env.DBConnectionParams)
	if err != nil {
		logger.Fatal(ctx, "failed to connect database", logger.Fields{
			logger.FieldError: err.Error(),
		})
	}

	setupServer(ctx, env, db)
}

func setupServer(
	ctx context.Context,
	envVars *config.EnvVariables,
	db *sql.DB,
) {
	app := fiber.New(fiber.Config{
		AppName: "SiteCrawler API",
	})

	// Middlewares
	app.Use(cors.New())
	app.Use(requestid.New())
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))

	// Inject RequestID into context for logging
	app.Use(func(c *fiber.Ctx) error {
		reqID := c.Locals("requestid")
		reqCtx := logger.WithCorrelationID(c.UserContext(), fmt.Sprintf("%v", reqID))
		c.SetUserContext(reqCtx)
		return c.Next()
	})

	// Register Routes
	// Initialize repositories
	crawlingSessionRepo := repository.CrawlingSessionRepository(postgres.NewCrawlingSessionRepo(db))
	pageRepo := repository.CrawlingSessionPageRepository(postgres.NewCrawlingSessionPageRepo(db))
	checkRepo := repository.CrawlingSessionCheckRepository(postgres.NewCrawlingSessionCheckRepo(db))
	auditRepo := repository.AuditCheckRepository(postgres.NewAuditRepo(db))
	viewRepo := repository.ViewRepository(postgres.NewViewRepo(db))
	statsRepo := repository.StatsRepository(postgres.NewStatsRepo(db))
	pageDetailsRepo := repository.PageDetailsRepository(postgres.NewPageDetailsRepo(db))

	sessionSvc := sessionsvc.NewService(
		sessionsvc.WithSessionRepository(crawlingSessionRepo),
		sessionsvc.WithPageRepository(pageRepo),
		sessionsvc.WithCheckRepository(checkRepo),
	)

	auditSvc := auditsvc.NewService(
		auditsvc.WithAuditCheckRepository(auditRepo),
	)

	viewSvc := viewsvc.NewService(
		viewsvc.WithViewRepository(viewRepo),
		viewsvc.WithPageRepository(pageRepo),
	)

	statsSvc := statssvc.NewService(
		statssvc.WithStatsRepository(statsRepo),
		statssvc.WithPageDetailsRepository(pageDetailsRepo),
	)
	metricsCtrl := stats.NewMetricsController()

	routeCtx := context.Background()
	routes.RegisterRoutes(routeCtx, app, metricsCtrl, auditSvc, sessionSvc, viewSvc, statsSvc)

	port := envVars.ServiceParams.Port
	logger.Info(ctx, "fiber server starting", logger.Fields{
		"port": port,
	})
	if err := app.Listen(":" + port); err != nil {
		logger.Fatal(ctx, "Failed to start server", logger.Fields{
			logger.FieldError: err.Error(),
		})
	}
}
