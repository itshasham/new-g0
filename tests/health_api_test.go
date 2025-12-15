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
	"sitecrawler/newgo/routes"
)

func TestHealthEndpoint(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		repoFactory    func() repository.HealthRepository
		expectedStatus int
		expectBody     bool
	}{
		{
			name: "healthy response",
			repoFactory: func() repository.HealthRepository {
				return repository.NewNoopHealthRepository()
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
		},
		{
			name: "repository failure",
			repoFactory: func() repository.HealthRepository {
				return repoFunc(func() error {
					return errors.New("db down")
				})
			},
			expectedStatus: http.StatusServiceUnavailable,
			expectBody:     false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := setupTestApp(tt.repoFactory())

			req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("fiber test request failed: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Fatalf("expected status %d got %d", tt.expectedStatus, resp.StatusCode)
			}

			if tt.expectBody {
				var payload dto.HealthResponse
				if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if payload.Status != "ok" {
					t.Fatalf("expected status ok got %s", payload.Status)
				}
			}
		})
	}
}

type repoFunc func() error

func (f repoFunc) Ping(_ context.Context) error {
	return f()
}

func setupTestApp(repo repository.HealthRepository) *fiber.App {
	app := fiber.New()
	service := services.NewHealthService(repo)
	controller := controllers.NewHealthController(service, nil)

	routes.Register(app, routes.Dependencies{
		Health: controller,
	})

	return app
}
