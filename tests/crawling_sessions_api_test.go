package tests

import (
sessionsDto "sitecrawler/newgo/dto/sessions"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"sitecrawler/newgo/controllers/health"
	"sitecrawler/newgo/controllers/sessions"
	"sitecrawler/newgo/internal/repository"
	sessionsvc "sitecrawler/newgo/internal/services/sessions"
	"sitecrawler/newgo/models"
	"sitecrawler/newgo/routes"
)

// =============================================================================
// CREATE CRAWLING SESSION TESTS
// =============================================================================

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
			app := setupCrawlingSessionApp(tt.repoFactory, tt.seed, nil, nil)

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
				var out sessionsDto.CrawlingSessionResponse
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

// =============================================================================
// GET CRAWLING SESSION TESTS
// =============================================================================

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
				var out sessionsDto.CrawlingSessionResponse
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
			app := setupCrawlingSessionApp(tt.repoFactory, tt.seed, nil, nil)
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

// =============================================================================
// LIST CRAWLING SESSION PAGES TESTS
// =============================================================================

func TestListCrawlingSessionPages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		path           string
		pageRepo       repository.CrawlingSessionPageRepository
		expectedStatus int
		assertBody     func(*testing.T, *http.Response)
	}{
		{
			name: "success",
			path: "/api/crawling_sessions/10/pages?page=1&page_limit=2",
			pageRepo: fakePageRepo{
				pages: []models.Page{
					{ID: 1, CrawlingSessionID: 10, URL: "https://example.com", ResponseCode: 200},
				},
				total: 1,
			},
			expectedStatus: http.StatusOK,
			assertBody: func(t *testing.T, resp *http.Response) {
				var out sessionsDto.CrawlingSessionPagesResponse
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
			pageRepo:       fakePageRepo{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid filters",
			path:           "/api/crawling_sessions/2/pages?filters=not-json",
			pageRepo:       fakePageRepo{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "repo error",
			path: "/api/crawling_sessions/2/pages",
			pageRepo: fakePageRepo{
				err: errors.New("failure"),
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := setupCrawlingSessionApp(nil, nil, tt.pageRepo, nil)
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

// =============================================================================
// LIST CRAWLING SESSION CHECKS TESTS
// =============================================================================

func TestListCrawlingSessionChecks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		path           string
		checkRepo      repository.CrawlingSessionCheckRepository
		expectedStatus int
		assertBody     func(*testing.T, *http.Response)
	}{
		{
			name: "success",
			path: "/api/crawling_sessions/5/checks_with_pages",
			checkRepo: fakeChecksRepo{
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
				var out sessionsDto.CrawlingSessionChecksResponse
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
			checkRepo:      fakeChecksRepo{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid comparison id",
			path:           "/api/crawling_sessions/5/checks_with_pages?comparison_crawling_session_id=abc",
			checkRepo:      fakeChecksRepo{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid filters",
			path:           "/api/crawling_sessions/5/checks_with_pages?filters=???",
			checkRepo:      fakeChecksRepo{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "repo error",
			path: "/api/crawling_sessions/5/checks_with_pages",
			checkRepo: fakeChecksRepo{
				err: errors.New("boom"),
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := setupCrawlingSessionApp(nil, nil, nil, tt.checkRepo)
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

// =============================================================================
// HELPER: UNIFIED TEST APP SETUP
// =============================================================================

func setupCrawlingSessionApp(
	sessionRepoFactory func() repository.CrawlingSessionRepository,
	sessionSeed func(*repository.InMemoryCrawlingSessionRepository),
	pageRepo repository.CrawlingSessionPageRepository,
	checkRepo repository.CrawlingSessionCheckRepository,
) *fiber.App {
	// Session repository
	sessionRepo := repository.CrawlingSessionRepository(nil)
	if sessionRepoFactory != nil {
		sessionRepo = sessionRepoFactory()
	}
	if sessionRepo == nil {
		sessionRepo = repository.NewInMemoryCrawlingSessionRepository()
	}
	if sessionSeed != nil {
		if mem, ok := sessionRepo.(*repository.InMemoryCrawlingSessionRepository); ok {
			sessionSeed(mem)
		}
	}

	// Page repository
	if pageRepo == nil {
		pageRepo = repository.NewNoopCrawlingSessionPageRepository()
	}

	// Check repository
	if checkRepo == nil {
		checkRepo = repository.NewNoopCrawlingSessionCheckRepository()
	}

	app := fiber.New()

	// Health
	healthController := health.NewController(nil)

	// Crawling session service using unified service
	sessionService := sessionsvc.NewService(sessionRepo, pageRepo, checkRepo)
	crawlingCreateController := sessions.NewCreateController(sessionService, nil)
	crawlingGetController := sessions.NewGetController(sessionService, nil)
	pagesController := sessions.NewPagesController(sessionService, nil)
	checksController := sessions.NewChecksController(sessionService, nil)

	routes.Register(app, routes.Dependencies{
		Health:                healthController,
		CrawlingSessionCreate: crawlingCreateController,
		CrawlingSessionGet:    crawlingGetController,
		CrawlingSessionPages:  pagesController,
		CrawlingSessionChecks: checksController,
	})

	return app
}

// =============================================================================
// MOCK REPOSITORIES
// =============================================================================

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

func (f failingCrawlingRepo) ClaimPending(ctx context.Context, queueID, limit int) ([]models.CrawlingSession, error) {
	return nil, nil
}

func (f failingCrawlingRepo) ClaimStalled(ctx context.Context, queueID int, excludeIDs []int64, limit int) ([]models.CrawlingSession, error) {
	return nil, nil
}

func (f failingCrawlingRepo) MarkDone(ctx context.Context, id int64, reason string) error {
	return nil
}

func (f failingCrawlingRepo) UpdateSiteInfo(ctx context.Context, id int64, info repository.SiteInfo) error {
	return nil
}

func (f failingCrawlingRepo) UpdateProgress(ctx context.Context, id int64, d repository.ProgressDelta) error {
	return nil
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
