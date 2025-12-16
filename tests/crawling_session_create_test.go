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

func TestGetCrawlingSession(t *testing.T) {
	t.Parallel()

	successSeed := func(repo *repository.InMemoryCrawlingSessionRepository) {
		_ = repo.Create(context.Background(), &models.CrawlingSession{
			SearchKeywordURLID: 77,
			URL:                "https://example.com",
			Status:             "done",
			Queue:              1,
		})
	}

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		repoFactory    func() repository.CrawlingSessionRepository
		seed           func(*repository.InMemoryCrawlingSessionRepository)
		assertBody     func(t *testing.T, resp *http.Response)
	}{
		{
			name:           "found",
			path:           "/api/crawling_sessions/1",
			expectedStatus: http.StatusOK,
			seed:           successSeed,
			assertBody: func(t *testing.T, resp *http.Response) {
				var out dto.CrawlingSessionResponse
				if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if out.Data.ID != 1 {
					t.Fatalf("expected id 1 got %d", out.Data.ID)
				}
				if out.Data.URL != "https://example.com" {
					t.Fatalf("unexpected url %s", out.Data.URL)
				}
			},
		},
		{
			name:           "invalid id",
			path:           "/api/crawling_sessions/abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "not found",
			path:           "/api/crawling_sessions/42",
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "repo error",
			path: "/api/crawling_sessions/99",
			repoFactory: func() repository.CrawlingSessionRepository {
				return failingCrawlingRepo{getErr: errors.New("db offline")}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := setupCrawlingApp(tt.repoFactory, tt.seed)
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("fiber request failed: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Fatalf("expected status %d got %d", tt.expectedStatus, resp.StatusCode)
			}

			if tt.assertBody != nil {
				tt.assertBody(t, resp)
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

	crawlingService := services.NewCrawlingSessionCreateService(repo)
	crawlingController := controllers.NewCrawlingSessionCreateController(crawlingService, nil)
	crawlingQueryService := services.NewCrawlingSessionGetService(repo)
	crawlingShowController := controllers.NewCrawlingSessionGetController(crawlingQueryService, nil)

	routes.Register(app, routes.Dependencies{
		Health:                healthController,
		CrawlingSessionCreate: crawlingController,
		CrawlingSessionGet:    crawlingShowController,
	})

	return app
}

type failingCrawlingRepo struct {
	preventErr error
	createErr  error
	getErr     error
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

func (f failingCrawlingRepo) GetByID(ctx context.Context, id int64) (*models.CrawlingSession, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return nil, repository.ErrCrawlingSessionNotFound
}
