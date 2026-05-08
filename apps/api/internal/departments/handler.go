package departments

import (
	"database/sql"
	"errors"
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
	router.Patch("/:id", h.update)
	router.Delete("/:id", h.delete)
}

func (h *Handler) list(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, response.APIError{
			Code:    "UNAUTHORIZED",
			Message: "unauthorized",
		})
	}

	items, err := h.service.List(c.UserContext(), user.OrganizationID)
	if err != nil {
		return writeError(c, err)
	}

	return response.Success(c, items)
}

func (h *Handler) create(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, response.APIError{
			Code:    "UNAUTHORIZED",
			Message: "unauthorized",
		})
	}

	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, response.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "invalid request body",
		})
	}

	dept, err := h.service.Create(c.UserContext(), user.OrganizationID, req)
	if err != nil {
		return writeError(c, err)
	}

	return response.Success(c, dept)
}

func (h *Handler) update(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, response.APIError{
			Code:    "UNAUTHORIZED",
			Message: "unauthorized",
		})
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, response.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "invalid request body",
		})
	}

	dept, err := h.service.Update(c.UserContext(), user.OrganizationID, c.Params("id"), req)
	if err != nil {
		return writeError(c, err)
	}

	return response.Success(c, dept)
}

func (h *Handler) delete(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, response.APIError{
			Code:    "UNAUTHORIZED",
			Message: "unauthorized",
		})
	}

	if err := h.service.Delete(c.UserContext(), user.OrganizationID, c.Params("id")); err != nil {
		return writeError(c, err)
	}

	return response.Success(c, fiber.Map{"deleted": true})
}

func writeError(c *fiber.Ctx, err error) error {
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
	case err != nil && strings.Contains(err.Error(), "required"):
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
