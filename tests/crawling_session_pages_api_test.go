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

func TestListCrawlingSessionPages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		path           string
		repo           repository.CrawlingSessionPageRepository
		expectedStatus int
		assertBody     func(*testing.T, *http.Response)
	}{
		{
			name: "success",
			path: "/api/crawling_sessions/10/pages?page=1&page_limit=2",
			repo: fakePageRepo{
				pages: []models.Page{
					{ID: 1, CrawlingSessionID: 10, URL: "https://example.com", ResponseCode: 200},
				},
				total: 1,
			},
			expectedStatus: http.StatusOK,
			assertBody: func(t *testing.T, resp *http.Response) {
				var out dto.CrawlingSessionPagesResponse
				if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if len(out.Data.Pages) != 1 {
					t.Fatalf("expected 1 page got %d", len(out.Data.Pages))
				}
				if out.Data.PagesTotal != 1 {
					t.Fatalf("expected total 1 got %d", out.Data.PagesTotal)
				}
			},
		},
		{
			name:           "invalid id",
			path:           "/api/crawling_sessions/foo/pages",
			repo:           fakePageRepo{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid filters",
			path:           "/api/crawling_sessions/2/pages?filters=not-json",
			repo:           fakePageRepo{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "repo error",
			path: "/api/crawling_sessions/2/pages",
			repo: fakePageRepo{
				err: errors.New("failure"),
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := setupPagesApp(tt.repo)
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

func setupPagesApp(repo repository.CrawlingSessionPageRepository) *fiber.App {
	app := fiber.New()
	healthService := services.NewHealthService(repository.NewNoopHealthRepository())
	healthController := controllers.NewHealthController(healthService, nil)

	pageService := services.NewCrawlingSessionPagesService(repo)
	pageController := controllers.NewCrawlingSessionPagesController(pageService, nil)

	routes.Register(app, routes.Dependencies{
		Health:               healthController,
		CrawlingSessionPages: pageController,
	})

	return app
}

type fakePageRepo struct {
	pages []models.Page
	total int
	err   error
}

func (f fakePageRepo) List(ctx context.Context, params repository.PageListParams) ([]models.Page, int, error) {
	if f.err != nil {
		return nil, 0, f.err
	}
	return f.pages, f.total, nil
}
