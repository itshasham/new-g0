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

func TestPageDetails(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		path           string
		repo           repository.PageDetailsRepository
		expectedStatus int
		assertBody     func(*testing.T, *http.Response)
	}{
		{
			name: "success 2xx",
			path: "/api/pages/2/page_details?search_keyword_url_id=5&limit=10",
			repo: fakePageDetailsRepo{
				page: &models.Page{ID: 2, ResponseCode: 200},
				images: []models.PageImage{
					{ID: 1, PageID: 2, URL: "https://example.com/img.png"},
				},
				broken: []models.Page{
					{ID: 3, URL: "https://example.com/broken", ResponseCode: 404},
				},
			},
			expectedStatus: http.StatusOK,
			assertBody: func(t *testing.T, resp *http.Response) {
				var out dto.PageDetailsResponse
				if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
					t.Fatalf("decode: %v", err)
				}
				if len(out.Data.PageImages) != 1 {
					t.Fatalf("expected images")
				}
				if len(out.Data.BrokenPages) != 1 {
					t.Fatalf("expected broken pages")
				}
				if len(out.Data.PagesLinkedToBrokenPage) != 0 {
					t.Fatalf("unexpected referrers")
				}
			},
		},
		{
			name: "success non 2xx",
			path: "/api/pages/9/page_details?search_keyword_url_id=5",
			repo: fakePageDetailsRepo{
				page: &models.Page{ID: 9, ResponseCode: 500},
				referrers: []models.Page{
					{ID: 7, URL: "https://example.com/ref"},
				},
			},
			expectedStatus: http.StatusOK,
			assertBody: func(t *testing.T, resp *http.Response) {
				var out dto.PageDetailsResponse
				if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
					t.Fatalf("decode: %v", err)
				}
				if len(out.Data.PagesLinkedToBrokenPage) != 1 {
					t.Fatalf("expected referrers")
				}
			},
		},
		{
			name:           "invalid id",
			path:           "/api/pages/foo/page_details?search_keyword_url_id=5",
			expectedStatus: http.StatusBadRequest,
			repo:           fakePageDetailsRepo{},
		},
		{
			name:           "missing sku",
			path:           "/api/pages/2/page_details",
			expectedStatus: http.StatusBadRequest,
			repo:           fakePageDetailsRepo{},
		},
		{
			name: "not found",
			path: "/api/pages/2/page_details?search_keyword_url_id=5",
			repo: fakePageDetailsRepo{
				err: repository.ErrPageNotFound,
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "repo failure",
			path: "/api/pages/2/page_details?search_keyword_url_id=5",
			repo: fakePageDetailsRepo{
				err: errors.New("db down"),
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := setupPageDetailsApp(tt.repo)
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

func setupPageDetailsApp(repo repository.PageDetailsRepository) *fiber.App {
	app := fiber.New()

	healthService := services.NewHealthService(repository.NewNoopHealthRepository())
	healthController := controllers.NewHealthController(healthService, nil)

	if repo == nil {
		repo = fakePageDetailsRepo{}
	}

	service := services.NewPageDetailsService(repo)
	controller := controllers.NewPageDetailsController(service, nil)

	routes.Register(app, routes.Dependencies{
		Health:      healthController,
		PageDetails: controller,
	})

	return app
}

type fakePageDetailsRepo struct {
	page      *models.Page
	images    []models.PageImage
	broken    []models.Page
	referrers []models.Page
	err       error
}

func (f fakePageDetailsRepo) GetPageByIDAndSKU(ctx context.Context, pageID, skuID int64) (*models.Page, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.page == nil {
		return nil, repository.ErrPageNotFound
	}
	return f.page, nil
}

func (f fakePageDetailsRepo) GetPageImages(ctx context.Context, pageID int64, limit int) ([]models.PageImage, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.images, nil
}

func (f fakePageDetailsRepo) GetBrokenTargetsFrom(ctx context.Context, pageID int64, limit int) ([]models.Page, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.broken, nil
}

func (f fakePageDetailsRepo) GetReferrersToBroken(ctx context.Context, pageID int64, limit int) ([]models.Page, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.referrers, nil
}
