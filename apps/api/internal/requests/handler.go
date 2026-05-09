package requests

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
	router.Post("/:id/submit", h.submit)
	router.Post("/:id/approval-actions", h.approvalAction)
	router.Post("/:id/comments", h.createComment)
	router.Get("/:id/comments", h.listComments)
}

func (h *Handler) list(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}

	canViewAll := middleware.HasAnyRole(user, "owner", "admin", "auditor", "finance", "director", "procurement", "manager")
	result, err := h.service.List(c.UserContext(), user.OrganizationID, canViewAll, user.ID, ListFilters{
		Page:         parseInt(c.Query("page"), 1),
		PageSize:     parseInt(c.Query("page_size"), 20),
		Status:       c.Query("status"),
		DepartmentID: c.Query("department_id"),
		RequesterID:  c.Query("requester_id"),
		FromDate:     c.Query("from_date"),
		ToDate:       c.Query("to_date"),
		Search:       c.Query("search"),
	})
	if err != nil {
		return writeRequestError(c, err)
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
		return badRequest(c, "invalid request body")
	}

	item, err := h.service.Create(c.UserContext(), user.OrganizationID, user.ID, req)
	if err != nil {
		return writeRequestError(c, err)
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
		return writeRequestError(c, err)
	}

	if !middleware.HasAnyRole(user, "owner", "admin", "auditor", "finance", "director", "procurement", "manager") && item.RequesterID != user.ID {
		return forbidden(c)
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
		return badRequest(c, "invalid request body")
	}

	item, err := h.service.Get(c.UserContext(), user.OrganizationID, c.Params("id"))
	if err != nil {
		return writeRequestError(c, err)
	}

	if !middleware.HasAnyRole(user, "owner", "admin") && item.RequesterID != user.ID {
		return forbidden(c)
	}

	updated, err := h.service.Update(c.UserContext(), user.OrganizationID, user.ID, c.Params("id"), req)
	if err != nil {
		return writeRequestError(c, err)
	}

	return response.Success(c, updated)
}

func (h *Handler) submit(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}

	item, err := h.service.Get(c.UserContext(), user.OrganizationID, c.Params("id"))
	if err != nil {
		return writeRequestError(c, err)
	}

	if !middleware.HasAnyRole(user, "owner", "admin") && item.RequesterID != user.ID {
		return forbidden(c)
	}

	result, err := h.service.Submit(c.UserContext(), user.OrganizationID, user.ID, c.Params("id"))
	if err != nil {
		return writeRequestError(c, err)
	}

	return response.Success(c, result)
}

func (h *Handler) approvalAction(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}

	var req ApprovalActionRequest
	if err := c.BodyParser(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	result, err := h.service.ActOnRequest(c.UserContext(), user.OrganizationID, ActorContext{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
		Roles:        user.Roles,
	}, c.Params("id"), req)
	if err != nil {
		return writeRequestError(c, err)
	}

	return response.Success(c, result)
}

func (h *Handler) createComment(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}

	var req struct {
		Body string `json:"body"`
	}
	if err := c.BodyParser(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	item, err := h.service.AddComment(c.UserContext(), user.OrganizationID, ActorContext{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
		Roles:        user.Roles,
	}, c.Params("id"), req.Body)
	if err != nil {
		return writeRequestError(c, err)
	}

	return response.Success(c, item)
}

func (h *Handler) listComments(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}

	items, err := h.service.ListComments(c.UserContext(), user.OrganizationID, c.Params("id"))
	if err != nil {
		return writeRequestError(c, err)
	}

	return response.SuccessList(c, items, fiber.Map{
		"page":      1,
		"page_size": len(items),
		"total":     len(items),
	})
}

func (h *Handler) PendingApprovals(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}

	items, err := h.service.ListPendingApprovals(c.UserContext(), user.OrganizationID, ActorContext{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
		Roles:        user.Roles,
	})
	if err != nil {
		return writeRequestError(c, err)
	}

	return response.SuccessList(c, items, fiber.Map{
		"page":      1,
		"page_size": len(items),
		"total":     len(items),
	})
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

func badRequest(c *fiber.Ctx, message string) error {
	return response.Error(c, fiber.StatusBadRequest, response.APIError{
		Code:    "VALIDATION_ERROR",
		Message: message,
	})
}

func writeRequestError(c *fiber.Ctx, err error) error {
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
	case errors.Is(err, ErrForbidden) || strings.Contains(strings.ToLower(err.Error()), "forbidden") || strings.Contains(strings.ToLower(err.Error()), "cannot act"):
		return response.Error(c, fiber.StatusForbidden, response.APIError{
			Code:    "FORBIDDEN",
			Message: "forbidden",
		})
	case strings.Contains(strings.ToLower(err.Error()), "required") || strings.Contains(strings.ToLower(err.Error()), "must") || strings.Contains(strings.ToLower(err.Error()), "invalid") || strings.Contains(strings.ToLower(err.Error()), "only requester"):
		return response.Error(c, fiber.StatusBadRequest, response.APIError{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
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
