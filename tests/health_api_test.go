package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/controllers/health"
	"sitecrawler/newgo/dto"
	"sitecrawler/newgo/routes"
)

func TestHealthEndpoint(t *testing.T) {
	t.Parallel()

	app := setupHealthTestApp()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("fiber test request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, resp.StatusCode)
	}

	var payload dto.HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Status != "ok" {
		t.Fatalf("expected status ok got %s", payload.Status)
	}
}

func setupHealthTestApp() *fiber.App {
	app := fiber.New()

	// Use new simplified health package
	controller := health.NewController(nil)

	routes.Register(app, routes.Dependencies{
		Health: controller,
	})

	return app
}
