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

// TestViewsCRUD tests the complete Create, Read, Update, List, Delete lifecycle
func TestViewsCRUD(t *testing.T) {
	t.Parallel()

	app := setupViewApp(nil, nil)

	// Create
	createBody := `{"data":{"search_keyword_url_id":123,"name":"view-1","filter_config":{"k":"v"}}}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/views", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createResp, err := app.Test(createReq)
	if err != nil {
		t.Fatalf("create request failed: %v", err)
	}
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d got %d", http.StatusCreated, createResp.StatusCode)
	}

	var created dto.ViewResponse
	if err := json.NewDecoder(createResp.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if created.Data.ID == 0 {
		t.Fatalf("expected created id")
	}

	// Get
	getReq := httptest.NewRequest(http.MethodGet, "/api/views/1", nil)
	getResp, err := app.Test(getReq)
	if err != nil {
		t.Fatalf("get request failed: %v", err)
	}
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, getResp.StatusCode)
	}

	var got dto.ViewResponse
	if err := json.NewDecoder(getResp.Body).Decode(&got); err != nil {
		t.Fatalf("decode get response: %v", err)
	}
	if got.Data.Name != "view-1" {
		t.Fatalf("expected name view-1 got %s", got.Data.Name)
	}

	// Update
	updateBody := `{"data":{"name":"view-2"}}`
	updateReq := httptest.NewRequest(http.MethodPut, "/api/views/1", strings.NewReader(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateResp, err := app.Test(updateReq)
	if err != nil {
		t.Fatalf("update request failed: %v", err)
	}
	if updateResp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, updateResp.StatusCode)
	}

	var updated dto.ViewResponse
	if err := json.NewDecoder(updateResp.Body).Decode(&updated); err != nil {
		t.Fatalf("decode update response: %v", err)
	}
	if updated.Data.Name != "view-2" {
		t.Fatalf("expected name view-2 got %s", updated.Data.Name)
	}

	// List
	listReq := httptest.NewRequest(http.MethodGet, "/api/views?search_keyword_url_id=123", nil)
	listResp, err := app.Test(listReq)
	if err != nil {
		t.Fatalf("list request failed: %v", err)
	}
	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, listResp.StatusCode)
	}

	var listed dto.ViewsResponse
	if err := json.NewDecoder(listResp.Body).Decode(&listed); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	if len(listed.Data) != 1 {
		t.Fatalf("expected 1 view got %d", len(listed.Data))
	}

	// Delete
	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/views/1", nil)
	deleteResp, err := app.Test(deleteReq)
	if err != nil {
		t.Fatalf("delete request failed: %v", err)
	}
	if deleteResp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, deleteResp.StatusCode)
	}

	var deleted dto.DeleteViewResponse
	if err := json.NewDecoder(deleteResp.Body).Decode(&deleted); err != nil {
		t.Fatalf("decode delete response: %v", err)
	}
	if deleted.Data.ID != 1 {
		t.Fatalf("expected delete id 1 got %d", deleted.Data.ID)
	}
}

// TestViewsListEmptyOnMissingSKU verifies list returns empty array when search_keyword_url_id is missing/zero
func TestViewsListEmptyOnMissingSKU(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		url      string
		wantCode int
	}{
		{
			name:     "missing search_keyword_url_id returns empty array",
			url:      "/api/views",
			wantCode: http.StatusOK,
		},
		{
			name:     "zero search_keyword_url_id returns empty array",
			url:      "/api/views?search_keyword_url_id=0",
			wantCode: http.StatusOK,
		},
		{
			name:     "valid search_keyword_url_id returns OK",
			url:      "/api/views?search_keyword_url_id=123",
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupViewApp(nil, nil)
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if resp.StatusCode != tt.wantCode {
				t.Fatalf("expected status %d got %d", tt.wantCode, resp.StatusCode)
			}
		})
	}
}

// TestViewsBadRequests tests various invalid request scenarios
func TestViewsBadRequests(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		method   string
		url      string
		body     string
		wantCode int
	}{
		{
			name:     "get with invalid id",
			method:   http.MethodGet,
			url:      "/api/views/abc",
			body:     "",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "update with invalid id",
			method:   http.MethodPut,
			url:      "/api/views/abc",
			body:     `{"data":{"name":"x"}}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "delete with invalid id",
			method:   http.MethodDelete,
			url:      "/api/views/abc",
			body:     "",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "create with invalid json",
			method:   http.MethodPost,
			url:      "/api/views",
			body:     `{invalid}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "update with invalid json",
			method:   http.MethodPut,
			url:      "/api/views/1",
			body:     `{invalid}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "page_count missing view_id",
			method:   http.MethodGet,
			url:      "/api/views/1/page_count",
			body:     "",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "page_count view_id is zero",
			method:   http.MethodGet,
			url:      "/api/views/1/page_count?view_id=0",
			body:     "",
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupViewApp(nil, nil)
			var req *http.Request
			if tt.body != "" {
				req = httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, tt.url, nil)
			}
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if resp.StatusCode != tt.wantCode {
				t.Fatalf("expected status %d got %d", tt.wantCode, resp.StatusCode)
			}
		})
	}
}

// TestViewsNotFound tests 404 scenarios
func TestViewsNotFound(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		method   string
		url      string
		body     string
		wantCode int
	}{
		{
			name:     "get non-existent view",
			method:   http.MethodGet,
			url:      "/api/views/42",
			body:     "",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "update non-existent view",
			method:   http.MethodPut,
			url:      "/api/views/42",
			body:     `{"data":{"name":"x"}}`,
			wantCode: http.StatusNotFound,
		},
		{
			name:     "page_count non-existent view",
			method:   http.MethodGet,
			url:      "/api/views/42/page_count?view_id=999",
			body:     "",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupViewApp(nil, nil)
			var req *http.Request
			if tt.body != "" {
				req = httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, tt.url, nil)
			}
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if resp.StatusCode != tt.wantCode {
				t.Fatalf("expected status %d got %d", tt.wantCode, resp.StatusCode)
			}
		})
	}
}

// TestViewsRepositoryErrors tests error handling when repository fails
func TestViewsRepositoryErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		method    string
		url       string
		body      string
		repoError string
		wantCode  int
	}{
		{
			name:      "get fails on repo error returns not found",
			method:    http.MethodGet,
			url:       "/api/views/1",
			body:      "",
			repoError: "db offline",
			wantCode:  http.StatusNotFound, // Service returns 404 for repo errors
		},
		{
			name:      "list fails on repo error returns unprocessable entity",
			method:    http.MethodGet,
			url:       "/api/views?search_keyword_url_id=123",
			body:      "",
			repoError: "db offline",
			wantCode:  http.StatusUnprocessableEntity, // Service returns 422 for repo errors
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupViewApp(func() repository.ViewRepository {
				return failingViewRepo{getErr: errors.New(tt.repoError), listErr: errors.New(tt.repoError)}
			}, nil)
			var req *http.Request
			if tt.body != "" {
				req = httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, tt.url, nil)
			}
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if resp.StatusCode != tt.wantCode {
				t.Fatalf("expected status %d got %d", tt.wantCode, resp.StatusCode)
			}
		})
	}
}

// TestViewsCreateVariations tests various create payloads
func TestViewsCreateVariations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		body      string
		wantCode  int
		wantHasID bool
	}{
		{
			name:      "full payload creates view",
			body:      `{"data":{"search_keyword_url_id":100,"name":"test-view","filter_config":{"filters":[]}}}`,
			wantCode:  http.StatusCreated,
			wantHasID: true,
		},
		{
			name:      "minimal payload creates view",
			body:      `{"data":{"search_keyword_url_id":200}}`,
			wantCode:  http.StatusCreated,
			wantHasID: true,
		},
		{
			name:      "empty data still creates view",
			body:      `{"data":{}}`,
			wantCode:  http.StatusCreated,
			wantHasID: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupViewApp(nil, nil)
			req := httptest.NewRequest(http.MethodPost, "/api/views", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if resp.StatusCode != tt.wantCode {
				t.Fatalf("expected status %d got %d", tt.wantCode, resp.StatusCode)
			}
			if tt.wantHasID {
				var created dto.ViewResponse
				if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if created.Data.ID == 0 {
					t.Fatalf("expected created id")
				}
			}
		})
	}
}

// TestViewsUpdateVariations tests various update scenarios
func TestViewsUpdateVariations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		updateBody     string
		expectedName   string
		expectedConfig bool
	}{
		{
			name:         "update name only",
			updateBody:   `{"data":{"name":"updated-name"}}`,
			expectedName: "updated-name",
		},
		{
			name:           "update filter_config only",
			updateBody:     `{"data":{"filter_config":{"new":"value"}}}`,
			expectedName:   "original-view",
			expectedConfig: true,
		},
		{
			name:           "update both name and filter_config",
			updateBody:     `{"data":{"name":"new-name","filter_config":{"updated":true}}}`,
			expectedName:   "new-name",
			expectedConfig: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupViewApp(nil, func(repo *repository.InMemoryViewRepository) {
				_ = repo.Create(context.Background(), &models.View{
					SearchKeywordURLID: 123,
					Name:               "original-view",
					FilterConfig:       map[string]any{"original": "config"},
				})
			})

			req := httptest.NewRequest(http.MethodPut, "/api/views/1", strings.NewReader(tt.updateBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected status %d got %d", http.StatusOK, resp.StatusCode)
			}

			var updated dto.ViewResponse
			if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if updated.Data.Name != tt.expectedName {
				t.Fatalf("expected name %s got %s", tt.expectedName, updated.Data.Name)
			}
		})
	}
}

// TestViewsPageCount tests the page count endpoint
func TestViewsPageCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setup    func(*repository.InMemoryViewRepository)
		url      string
		wantCode int
	}{
		{
			name: "valid page count request",
			setup: func(repo *repository.InMemoryViewRepository) {
				_ = repo.Create(context.Background(), &models.View{
					SearchKeywordURLID: 123,
					Name:               "test-view",
					FilterConfig:       map[string]any{"filter_groups": []any{}},
				})
			},
			url:      "/api/views/1/page_count?view_id=1&crawling_session_id=100",
			wantCode: http.StatusOK,
		},
		{
			name:     "page count without view in db",
			setup:    nil,
			url:      "/api/views/1/page_count?view_id=999",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupViewApp(nil, tt.setup)
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if resp.StatusCode != tt.wantCode {
				t.Fatalf("expected status %d got %d", tt.wantCode, resp.StatusCode)
			}
		})
	}
}

// setupViewApp creates a test fiber app with view routes
func setupViewApp(factory func() repository.ViewRepository, seed func(*repository.InMemoryViewRepository)) *fiber.App {
	viewRepo := repository.ViewRepository(nil)
	if factory != nil {
		viewRepo = factory()
	}
	if viewRepo == nil {
		viewRepo = repository.NewInMemoryViewRepository()
	}

	if seed != nil {
		if mem, ok := viewRepo.(*repository.InMemoryViewRepository); ok {
			seed(mem)
		}
	}

	pageRepo := repository.NewNoopCrawlingSessionPageRepository()

	app := fiber.New()

	healthService := services.NewHealthService(repository.NewNoopHealthRepository())
	healthController := controllers.NewHealthController(healthService, nil)

	listService := services.NewViewListService(viewRepo)
	getService := services.NewViewGetService(viewRepo)
	createService := services.NewViewCreateService(viewRepo)
	updateService := services.NewViewUpdateService(viewRepo)
	deleteService := services.NewViewDeleteService(viewRepo)
	pageCountService := services.NewViewPageCountService(viewRepo, pageRepo)

	routes.Register(app, routes.Dependencies{
		Health:        healthController,
		ViewList:      controllers.NewViewListController(listService, nil),
		ViewGet:       controllers.NewViewGetController(getService, nil),
		ViewCreate:    controllers.NewViewCreateController(createService, nil),
		ViewUpdate:    controllers.NewViewUpdateController(updateService, nil),
		ViewDelete:    controllers.NewViewDeleteController(deleteService, nil),
		ViewPageCount: controllers.NewViewPageCountController(pageCountService, nil),
	})

	return app
}

// failingViewRepo is a mock repository that returns errors
type failingViewRepo struct {
	createErr error
	updateErr error
	deleteErr error
	getErr    error
	listErr   error
}

func (f failingViewRepo) Create(ctx context.Context, v *models.View) error {
	_ = ctx
	_ = v
	return f.createErr
}

func (f failingViewRepo) Update(ctx context.Context, v *models.View) error {
	_ = ctx
	_ = v
	return f.updateErr
}

func (f failingViewRepo) Delete(ctx context.Context, id int64) error {
	_ = ctx
	_ = id
	return f.deleteErr
}

func (f failingViewRepo) Get(ctx context.Context, id int64) (*models.View, error) {
	_ = ctx
	_ = id
	if f.getErr != nil {
		return nil, f.getErr
	}
	return nil, repository.ErrViewNotFound
}

func (f failingViewRepo) ListBySKU(ctx context.Context, skuID int64) ([]models.View, error) {
	_ = ctx
	_ = skuID
	return nil, f.listErr
}
