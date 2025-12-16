package tests

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/controllers"
	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/internal/repository"
	"sitecrawler/newgo/internal/services"
	"sitecrawler/newgo/models"
	"sitecrawler/newgo/routes"
)

func TestListCrawlingSessionChecks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		path           string
		repo           repository.CrawlingSessionCheckRepository
		expectedStatus int
		assertBody     func(*testing.T, *http.Response)
	}{
		{
			name: "success",
			path: "/api/crawling_sessions/5/checks_with_pages",
			repo: fakeChecksRepo{
				checks: []models.CheckWithPages{
					{
						ID:   1,
						Name: "Broken Links",
						Pages: []models.Page{
							{ID: 100, CrawlingSessionID: 5, URL: "https://example.com", ResponseCode: 404},
						},
					},
				},
			},
			expectedStatus: http.StatusOK,
			assertBody: func(t *testing.T, resp *http.Response) {
				var out dto.CrawlingSessionChecksResponse
				if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if len(out.Data.Checks) != 1 {
					t.Fatalf("expected 1 check got %d", len(out.Data.Checks))
				}
				if len(out.Data.Checks[0].Pages) != 1 {
					t.Fatalf("expected pages for check")
				}
			},
		},
		{
			name:           "invalid id",
			path:           "/api/crawling_sessions/foo/checks_with_pages",
			repo:           fakeChecksRepo{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid comparison id",
			path:           "/api/crawling_sessions/5/checks_with_pages?comparison_crawling_session_id=abc",
			repo:           fakeChecksRepo{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid filters",
			path:           "/api/crawling_sessions/5/checks_with_pages?filters=???",
			repo:           fakeChecksRepo{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "repo error",
			path: "/api/crawling_sessions/5/checks_with_pages",
			repo: fakeChecksRepo{
				err: errors.New("boom"),
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := setupChecksApp(tt.repo)
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

func setupChecksApp(repo repository.CrawlingSessionCheckRepository) *fiber.App {
	app := fiber.New()
	healthService := services.NewHealthService(repository.NewNoopHealthRepository())
	healthController := controllers.NewHealthController(healthService, nil)

	checkService := services.NewCrawlingSessionChecksService(repo)
	checkController := controllers.NewCrawlingSessionChecksController(checkService, nil)

	routes.Register(app, routes.Dependencies{
		Health:                healthController,
		CrawlingSessionChecks: checkController,
	})

	return app
}

type fakeChecksRepo struct {
	checks []models.CheckWithPages
	err    error
}

func (f fakeChecksRepo) ChecksWithPages(ctx context.Context, params repository.ChecksWithPagesParams) ([]models.CheckWithPages, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.checks, nil
}
