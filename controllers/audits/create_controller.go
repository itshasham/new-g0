package audits

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	auditsDto "sitecrawler/newgo/dto/audits"
	"sitecrawler/newgo/internal/services/audits"

	"github.com/gofiber/fiber/v2"
)

type CreateController struct {
	service audits.Service
	logger  *slog.Logger
}

func NewCreateController(service audits.Service, logger *slog.Logger) *CreateController {
	if service == nil {
		panic("audit check create service required")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &CreateController{service: service, logger: logger}
}

func (c *CreateController) Create(ctx *fiber.Ctx) error {
	var request auditsDto.CreateAuditCheckRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid json payload"})
	}

	if err := validateCreateAuditCheckRequest(request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := c.service.Create(ctx.Context(), request)
	if err != nil {
		c.logger.Error("audit check create failed", "error", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if resp.Body == nil {
		return ctx.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Message})
	}
	return ctx.Status(resp.StatusCode).JSON(resp.Body)
}

func validateCreateAuditCheckRequest(req auditsDto.CreateAuditCheckRequest) error {
	if req.Data.SearchKeywordURLID == 0 {
		return errors.New("search_keyword_url_id is required")
	}
	if strings.TrimSpace(req.Data.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(req.Data.Category) == "" {
		return errors.New("category is required")
	}
	return nil
}
