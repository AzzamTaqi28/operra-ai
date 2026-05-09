package audit

import (
	"database/sql"
	"errors"
	"strconv"

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
}

func (h *Handler) list(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}
	if !middleware.HasAnyRole(user, "owner", "admin", "finance", "auditor") {
		return forbidden(c)
	}

	result, err := h.service.List(c.UserContext(), user.OrganizationID, ListFilters{
		Page:        parseInt(c.Query("page"), 1),
		PageSize:    parseInt(c.Query("page_size"), 20),
		Action:      c.Query("action"),
		EntityType:  c.Query("entity_type"),
		EntityID:    c.Query("entity_id"),
		FromDate:    c.Query("from_date"),
		ToDate:      c.Query("to_date"),
		ActorUserID: c.Query("actor_user_id"),
	})
	if err != nil {
		return writeError(c, err)
	}

	return response.SuccessList(c, result.Items, result.Pagination)
}

func unauthorized(c *fiber.Ctx) error {
	return response.Error(c, fiber.StatusUnauthorized, response.APIError{
		Code:    "UNAUTHORIZED",
		Message: "unauthorized",
	})
}

func forbidden(c *fiber.Ctx) error {
	return response.Error(c, fiber.StatusForbidden, response.APIError{
		Code:    "FORBIDDEN",
		Message: "forbidden",
	})
}

func writeError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return response.Error(c, fiber.StatusNotFound, response.APIError{
			Code:    "NOT_FOUND",
			Message: "resource not found",
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
