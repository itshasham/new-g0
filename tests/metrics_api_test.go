package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/controllers"
)

func TestMetricsEndpoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		checkBody      func(*testing.T, string)
	}{
		{
			name:           "metrics endpoint returns prometheus format",
			method:         http.MethodGet,
			path:           "/metrics",
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				// Prometheus metrics should contain standard Go metrics
				if !strings.Contains(body, "go_") {
					t.Error("expected go_ metrics prefix in response")
				}
				// Should have HELP and TYPE annotations
				if !strings.Contains(body, "# HELP") {
					t.Error("expected # HELP annotation in prometheus format")
				}
				if !strings.Contains(body, "# TYPE") {
					t.Error("expected # TYPE annotation in prometheus format")
				}
			},
		},
		{
			name:           "metrics endpoint with query params still works",
			method:         http.MethodGet,
			path:           "/metrics?format=text",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupMetricsApp()
			req := httptest.NewRequest(tt.method, tt.path, nil)

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Fatalf("expected status %d got %d", tt.expectedStatus, resp.StatusCode)
			}

			if tt.checkBody != nil {
				body, _ := io.ReadAll(resp.Body)
				tt.checkBody(t, string(body))
			}
		})
	}
}

func TestMetricsContentType(t *testing.T) {
	t.Parallel()

	app := setupMetricsApp()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	contentType := resp.Header.Get("Content-Type")
	// Prometheus metrics can return various content types
	// Usually "text/plain" or "text/plain; version=0.0.4; charset=utf-8"
	if !strings.Contains(contentType, "text/plain") {
		t.Fatalf("expected content-type to contain text/plain, got %s", contentType)
	}
}

func setupMetricsApp() *fiber.App {
	app := fiber.New()

	metricsCtrl := controllers.NewMetricsController()

	// Register metrics route directly without going through full deps
	app.Get("/metrics", metricsCtrl.Metrics)

	return app
}
