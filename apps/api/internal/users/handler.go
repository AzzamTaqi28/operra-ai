package users

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"operra/api/internal/platform/middleware"
	"operra/api/internal/platform/response"
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
	router.Patch("/:id", h.update)
	router.Post("/:id/roles", h.assignRoles)
	router.Delete("/:id/roles/:role_key", h.removeRole)
}

func (h *Handler) list(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}

	filters := ListFilters{
		Page:         parseInt(c.Query("page"), 1),
		PageSize:     parseInt(c.Query("page_size"), 20),
		Status:       c.Query("status"),
		Role:         c.Query("role"),
		DepartmentID: c.Query("department_id"),
		Search:       c.Query("search"),
	}

	result, err := h.service.List(c.UserContext(), user.OrganizationID, filters)
	if err != nil {
		return writeUserError(c, err)
	}

	return response.SuccessList(c, result.Items, result.Pagination)
}

func (h *Handler) create(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}

	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, response.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "invalid request body",
		})
	}

	item, err := h.service.Create(c.UserContext(), user.OrganizationID, req)
	if err != nil {
		return writeUserError(c, err)
	}

	return response.Success(c, item)
}

func (h *Handler) get(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}

	item, err := h.service.Get(c.UserContext(), user.OrganizationID, c.Params("id"))
	if err != nil {
		return writeUserError(c, err)
	}

	return response.Success(c, item)
}

func (h *Handler) update(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, response.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "invalid request body",
		})
	}

	item, err := h.service.Update(c.UserContext(), user.OrganizationID, c.Params("id"), req)
	if err != nil {
		return writeUserError(c, err)
	}

	return response.Success(c, item)
}

func (h *Handler) assignRoles(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}

	var req struct {
		RoleKeys []string `json:"role_keys"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, response.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "invalid request body",
		})
	}

	if err := h.service.AssignRoles(c.UserContext(), user.OrganizationID, c.Params("id"), req.RoleKeys); err != nil {
		return writeUserError(c, err)
	}

	return response.Success(c, fiber.Map{"updated": true})
}

func (h *Handler) removeRole(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}

	if err := h.service.RemoveRole(c.UserContext(), user.OrganizationID, c.Params("id"), c.Params("role_key")); err != nil {
		return writeUserError(c, err)
	}

	return response.Success(c, fiber.Map{"deleted": true})
}

func unauthorized(c *fiber.Ctx) error {
	return response.Error(c, fiber.StatusUnauthorized, response.APIError{
		Code:    "UNAUTHORIZED",
		Message: "unauthorized",
	})
}

func writeUserError(c *fiber.Ctx, err error) error {
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
	case strings.Contains(strings.ToLower(err.Error()), "required") || strings.Contains(strings.ToLower(err.Error()), "invalid") || strings.Contains(strings.ToLower(err.Error()), "at least one field"):
		return response.Error(c, fiber.StatusBadRequest, response.APIError{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
	default:
		return response.Error(c, fiber.StatusInternalServerError, response.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "internal server error",
		})
	}
}

func parseInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
