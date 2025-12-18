package tests

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	statsDto "sitecrawler/newgo/controllers/dto/stats"
	"testing"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/internal/repository"
	statssvc "sitecrawler/newgo/internal/services/stats"
	"sitecrawler/newgo/routes"
)

func TestStatsAPI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		path           string
		repo           repository.StatsRepository
		expectedStatus int
		assertBody     func(*testing.T, *http.Response)
	}{
		{
			name: "success",
			path: "/api/stats?crawling_session_id=1",
			repo: fakeStatsRepo{
				result: map[string]any{"pages": 10},
			},
			expectedStatus: http.StatusOK,
			assertBody: func(t *testing.T, resp *http.Response) {
				var out statsDto.StatsResponse
				if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				val, ok := out.Data["pages"].(float64)
				if !ok || int(val) != 10 {
					t.Fatalf("expected pages 10 got %v", out.Data["pages"])
				}
			},
		},
		{
			name:           "missing session id",
			path:           "/api/stats",
			repo:           fakeStatsRepo{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid filters",
			path:           "/api/stats?crawling_session_id=1&filters=notjson",
			repo:           fakeStatsRepo{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid prefilters",
			path:           "/api/stats?crawling_session_id=1&prefilters=??",
			repo:           fakeStatsRepo{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid comparison id",
			path:           "/api/stats?crawling_session_id=1&comparison_crawling_session_id=abc",
			repo:           fakeStatsRepo{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "repo error",
			path: "/api/stats?crawling_session_id=1",
			repo: fakeStatsRepo{
				err: errors.New("query failed"),
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := setupStatsApp(tt.repo)
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

func setupStatsApp(repo repository.StatsRepository) *fiber.App {
	app := fiber.New()

	if repo == nil {
		repo = fakeStatsRepo{}
	}

	// Use unified stats service
	statsService := statssvc.NewService(
		statssvc.WithStatsRepository(repo),
		statssvc.WithPageDetailsRepository(repository.NewNoopPageDetailsRepository()),
	)
	routes.RegisterRoutes(context.Background(), app, nil, nil, nil, nil, statsService)

	return app
}

type fakeStatsRepo struct {
	result map[string]any
	err    error
}

func (f fakeStatsRepo) Fetch(ctx context.Context, params repository.StatsQueryParams) (map[string]any, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.result == nil {
		return map[string]any{}, nil
	}
	return f.result, nil
}
