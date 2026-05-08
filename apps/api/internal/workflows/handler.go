package workflows

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"operra/api/internal/platform/middleware"
	"operra/api/internal/platform/response"
	configworkflow "operra/api/internal/workflow"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Get("", h.list)
	router.Post("", h.create)
	router.Get("/:id", h.get)
	router.Get("/:id/versions", h.versions)
	router.Post("/:id/versions", h.createVersion)
	router.Post("/:id/versions/:version_id/activate", h.activate)
	router.Post("/validate", h.validate)
}

func (h *Handler) list(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}

	items, err := h.service.List(c.UserContext(), user.OrganizationID, c.Query("type"), c.Query("status"))
	if err != nil {
		return writeWorkflowError(c, err)
	}

	return response.SuccessList(c, items, fiber.Map{
		"page":      1,
		"page_size": len(items),
		"total":     len(items),
	})
}

func (h *Handler) get(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}

	item, err := h.service.Get(c.UserContext(), user.OrganizationID, c.Params("id"))
	if err != nil {
		return writeWorkflowError(c, err)
	}

	return response.Success(c, item)
}

func (h *Handler) versions(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}

	items, err := h.service.ListVersions(c.UserContext(), user.OrganizationID, c.Params("id"))
	if err != nil {
		return writeWorkflowError(c, err)
	}

	return response.SuccessList(c, items, fiber.Map{
		"page":      1,
		"page_size": len(items),
		"total":     len(items),
	})
}

func (h *Handler) create(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}

	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	result, err := h.service.CreateWorkflow(c.UserContext(), user.OrganizationID, user.ID, req)
	if err != nil {
		return writeWorkflowError(c, err)
	}

	return response.Success(c, result)
}

func (h *Handler) createVersion(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}

	var req CreateVersionRequest
	if err := c.BodyParser(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	result, err := h.service.CreateVersion(c.UserContext(), user.OrganizationID, user.ID, c.Params("id"), req)
	if err != nil {
		return writeWorkflowError(c, err)
	}

	return response.Success(c, result)
}

func (h *Handler) activate(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}

	if err := h.service.ActivateVersion(c.UserContext(), user.OrganizationID, c.Params("id"), c.Params("version_id")); err != nil {
		return writeWorkflowError(c, err)
	}

	return response.Success(c, fiber.Map{"activated": true})
}

func (h *Handler) validate(c *fiber.Ctx) error {
	var req struct {
		Type       string          `json:"type"`
		ConfigJSON json.RawMessage `json:"config_json"`
	}
	if err := c.BodyParser(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	var cfg configworkflow.Config
	if len(req.ConfigJSON) > 0 {
		if err := json.Unmarshal(req.ConfigJSON, &cfg); err != nil {
			return badRequest(c, "config_json must be valid JSON")
		}
	}
	if req.Type != "" {
		cfg.Type = req.Type
	}

	result := h.service.validate(cfg)
	return response.Success(c, fiber.Map{
		"validation":      result,
		"mermaid_diagram": h.service.generate(cfg),
	})
}

func unauthorized(c *fiber.Ctx) error {
	return response.Error(c, fiber.StatusUnauthorized, response.APIError{
		Code:    "UNAUTHORIZED",
		Message: "unauthorized",
	})
}

func badRequest(c *fiber.Ctx, message string) error {
	return response.Error(c, fiber.StatusBadRequest, response.APIError{
		Code:    "VALIDATION_ERROR",
		Message: message,
	})
}

func writeWorkflowError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return response.Error(c, fiber.StatusNotFound, response.APIError{
			Code:    "NOT_FOUND",
			Message: "resource not found",
		})
	case errors.Is(err, ErrConflict) || strings.Contains(strings.ToLower(err.Error()), "conflict"):
		return response.Error(c, fiber.StatusConflict, response.APIError{
			Code:    "CONFLICT",
			Message: "resource already exists",
		})
	case strings.Contains(strings.ToLower(err.Error()), "workflow invalid"):
		var invalid ErrInvalidWorkflow
		if errors.As(err, &invalid) {
			return response.Error(c, fiber.StatusUnprocessableEntity, response.APIError{
				Code:    "WORKFLOW_INVALID",
				Message: "workflow invalid",
				Details: invalid.Validation.Errors,
			})
		}
		return response.Error(c, fiber.StatusUnprocessableEntity, response.APIError{
			Code:    "WORKFLOW_INVALID",
			Message: "workflow invalid",
		})
	default:
		return response.Error(c, fiber.StatusInternalServerError, response.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "internal server error",
		})
	}
}
