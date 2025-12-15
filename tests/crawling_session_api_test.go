package tests

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/controllers"
	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"
	"sitecrawler/newgo/internal/services"
	"sitecrawler/newgo/models"
	"sitecrawler/newgo/routes"
)

func TestCreateCrawlingSession(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		body           string
		expectedStatus int
		expectBody     bool
		repoFactory    func() repository.CrawlingSessionRepository
		seed           func(*repository.InMemoryCrawlingSessionRepository)
	}{
		{
			name:           "successful creation",
			body:           `{"data":{"search_keyword_url_id":123,"url":"https://example.com","queue":1,"options":{"depth":2}}}`,
			expectedStatus: http.StatusCreated,
			expectBody:     true,
		},
		{
			name:           "validation failure",
			body:           `{"data":{"search_keyword_url_id":0,"url":"","queue":1}}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "session already in progress",
			body: `{"data":{"search_keyword_url_id":777,"url":"https://example.com"}}`,
			seed: func(repo *repository.InMemoryCrawlingSessionRepository) {
				_ = repo.PreventInProgress(context.Background(), 777)
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "repository prevent failure",
			body: `{"data":{"search_keyword_url_id":999,"url":"https://example.com"}}`,
			repoFactory: func() repository.CrawlingSessionRepository {
				return failingCrawlingRepo{
					preventErr: errors.New("in progress"),
				}
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "repository create failure",
			body: `{"data":{"search_keyword_url_id":555,"url":"https://example.com"}}`,
			repoFactory: func() repository.CrawlingSessionRepository {
				return failingCrawlingRepo{
					createErr: errors.New("db write failed"),
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := setupCrawlingApp(tt.repoFactory, tt.seed)

			req := httptest.NewRequest(http.MethodPost, "/api/crawling_sessions", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("fiber request failed: %v", err)
			}
			if resp.StatusCode != tt.expectedStatus {
				t.Fatalf("expected status %d got %d", tt.expectedStatus, resp.StatusCode)
			}

			if tt.expectBody {
				var out dto.CrawlingSessionResponse
				if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if out.Data.SearchKeywordURLID != 123 {
					t.Fatalf("expected sku 123 got %d", out.Data.SearchKeywordURLID)
				}
				if out.Data.Status != "pending" {
					t.Fatalf("expected status pending got %s", out.Data.Status)
				}
			}
		})
	}
}

func setupCrawlingApp(factory func() repository.CrawlingSessionRepository, seed func(*repository.InMemoryCrawlingSessionRepository)) *fiber.App {
	repo := repository.CrawlingSessionRepository(nil)
	if factory != nil {
		repo = factory()
	}
	if repo == nil {
		repo = repository.NewInMemoryCrawlingSessionRepository()
	}

	if seed != nil {
		if mem, ok := repo.(*repository.InMemoryCrawlingSessionRepository); ok {
			seed(mem)
		}
	}

	app := fiber.New()

	healthService := services.NewHealthService(repository.NewNoopHealthRepository())
	healthController := controllers.NewHealthController(healthService, nil)

	crawlingService := services.NewCrawlingSessionService(repo)
	crawlingController := controllers.NewCrawlingSessionController(crawlingService, nil)

	routes.Register(app, routes.Dependencies{
		Health:           healthController,
		CrawlingSessions: crawlingController,
	})

	return app
}

type failingCrawlingRepo struct {
	preventErr error
	createErr  error
}

func (f failingCrawlingRepo) PreventInProgress(ctx context.Context, skuID int64) error {
	return f.preventErr
}

func (f failingCrawlingRepo) Create(ctx context.Context, session *models.CrawlingSession) error {
	if f.createErr != nil {
		return f.createErr
	}
	return errors.New("create not supported")
}
