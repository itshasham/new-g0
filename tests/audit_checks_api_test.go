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

func TestAuditChecksCRUD(t *testing.T) {
	t.Parallel()

	app := setupAuditApp(nil, nil)

	createBody := `{"data":{"search_keyword_url_id":123,"name":"check-1","category":"cat-1","filter_config":{"k":"v"}}}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/audit_checks", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createResp, err := app.Test(createReq)
	if err != nil {
		t.Fatalf("create request failed: %v", err)
	}
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d got %d", http.StatusCreated, createResp.StatusCode)
	}

	var created dto.AuditCheckResponse
	if err := json.NewDecoder(createResp.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if created.Data.ID == 0 {
		t.Fatalf("expected created id")
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/audit_checks/1", nil)
	getResp, err := app.Test(getReq)
	if err != nil {
		t.Fatalf("get request failed: %v", err)
	}
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, getResp.StatusCode)
	}

	var got dto.AuditCheckResponse
	if err := json.NewDecoder(getResp.Body).Decode(&got); err != nil {
		t.Fatalf("decode get response: %v", err)
	}
	if got.Data.Name != "check-1" {
		t.Fatalf("expected name check-1 got %s", got.Data.Name)
	}

	updateBody := `{"data":{"name":"check-2"}}`
	updateReq := httptest.NewRequest(http.MethodPut, "/api/audit_checks/1", strings.NewReader(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateResp, err := app.Test(updateReq)
	if err != nil {
		t.Fatalf("update request failed: %v", err)
	}
	if updateResp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, updateResp.StatusCode)
	}

	var updated dto.AuditCheckResponse
	if err := json.NewDecoder(updateResp.Body).Decode(&updated); err != nil {
		t.Fatalf("decode update response: %v", err)
	}
	if updated.Data.Name != "check-2" {
		t.Fatalf("expected name check-2 got %s", updated.Data.Name)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/audit_checks?search_keyword_url_id=123", nil)
	listResp, err := app.Test(listReq)
	if err != nil {
		t.Fatalf("list request failed: %v", err)
	}
	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, listResp.StatusCode)
	}

	var listed dto.AuditChecksResponse
	if err := json.NewDecoder(listResp.Body).Decode(&listed); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	if len(listed.Data) != 1 {
		t.Fatalf("expected 1 audit check got %d", len(listed.Data))
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/audit_checks/1", nil)
	deleteResp, err := app.Test(deleteReq)
	if err != nil {
		t.Fatalf("delete request failed: %v", err)
	}
	if deleteResp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, deleteResp.StatusCode)
	}

	var deleted dto.DeleteAuditCheckResponse
	if err := json.NewDecoder(deleteResp.Body).Decode(&deleted); err != nil {
		t.Fatalf("decode delete response: %v", err)
	}
	if deleted.Data.ID != 1 {
		t.Fatalf("expected delete id 1 got %d", deleted.Data.ID)
	}
}

func TestAuditChecksBadRequests(t *testing.T) {
	t.Parallel()

	app := setupAuditApp(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/audit_checks", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("fiber request failed: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, resp.StatusCode)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/audit_checks/abc", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("fiber request failed: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, resp.StatusCode)
	}

	req = httptest.NewRequest(http.MethodPut, "/api/audit_checks/abc", strings.NewReader(`{"data":{"name":"x"}}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("fiber request failed: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, resp.StatusCode)
	}

	req = httptest.NewRequest(http.MethodDelete, "/api/audit_checks/abc", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("fiber request failed: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, resp.StatusCode)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/audit_checks", strings.NewReader(`{"data":{"search_keyword_url_id":0}}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("fiber request failed: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestAuditChecksNotFound(t *testing.T) {
	t.Parallel()

	app := setupAuditApp(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/audit_checks/42", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("fiber request failed: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status %d got %d", http.StatusNotFound, resp.StatusCode)
	}

	req = httptest.NewRequest(http.MethodPut, "/api/audit_checks/42", strings.NewReader(`{"data":{"name":"x"}}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("fiber request failed: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status %d got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestAuditChecksRepositoryErrors(t *testing.T) {
	t.Parallel()

	app := setupAuditApp(func() repository.AuditCheckRepository {
		return failingAuditRepo{getErr: errors.New("db offline")}
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/audit_checks/1", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("fiber request failed: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status %d got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}

func setupAuditApp(factory func() repository.AuditCheckRepository, seed func(*repository.InMemoryAuditCheckRepository)) *fiber.App {
	repo := repository.AuditCheckRepository(nil)
	if factory != nil {
		repo = factory()
	}
	if repo == nil {
		repo = repository.NewInMemoryAuditCheckRepository()
	}

	if seed != nil {
		if mem, ok := repo.(*repository.InMemoryAuditCheckRepository); ok {
			seed(mem)
		}
	}

	app := fiber.New()

	healthService := services.NewHealthService(repository.NewNoopHealthRepository())
	healthController := controllers.NewHealthController(healthService, nil)

	listService := services.NewAuditCheckListService(repo)
	getService := services.NewAuditCheckGetService(repo)
	createService := services.NewAuditCheckCreateService(repo)
	updateService := services.NewAuditCheckUpdateService(repo)
	deleteService := services.NewAuditCheckDeleteService(repo)

	routes.Register(app, routes.Dependencies{
		Health:           healthController,
		AuditCheckList:   controllers.NewAuditCheckListController(listService, nil),
		AuditCheckGet:    controllers.NewAuditCheckGetController(getService, nil),
		AuditCheckCreate: controllers.NewAuditCheckCreateController(createService, nil),
		AuditCheckUpdate: controllers.NewAuditCheckUpdateController(updateService, nil),
		AuditCheckDelete: controllers.NewAuditCheckDeleteController(deleteService, nil),
	})

	return app
}

type failingAuditRepo struct {
	createErr error
	updateErr error
	deleteErr error
	getErr    error
	listErr   error
}

func (f failingAuditRepo) Create(ctx context.Context, ac *models.AuditCheck) error {
	_ = ctx
	_ = ac
	return f.createErr
}

func (f failingAuditRepo) Update(ctx context.Context, ac *models.AuditCheck) error {
	_ = ctx
	_ = ac
	return f.updateErr
}

func (f failingAuditRepo) Delete(ctx context.Context, id int64) error {
	_ = ctx
	_ = id
	return f.deleteErr
}

func (f failingAuditRepo) Get(ctx context.Context, id int64) (*models.AuditCheck, error) {
	_ = ctx
	_ = id
	if f.getErr != nil {
		return nil, f.getErr
	}
	return nil, repository.ErrAuditCheckNotFound
}

func (f failingAuditRepo) ListBySKU(ctx context.Context, skuID int64) ([]models.AuditCheck, error) {
	_ = ctx
	_ = skuID
	return nil, f.listErr
}

func (f failingAuditRepo) ListByIDsAndSKUs(ctx context.Context, ids, skuIDs []int64) ([]models.AuditCheck, error) {
	_ = ctx
	_ = ids
	_ = skuIDs
	return nil, nil
}
