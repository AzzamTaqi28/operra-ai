package exports

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
	router.Get("/purchase-requests.csv", h.purchaseRequests)
	router.Get("/approval-history.csv", h.approvalHistory)
	router.Get("/audit-logs.csv", h.auditLogs)
}

func (h *Handler) purchaseRequests(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}
	if !middleware.HasAnyRole(user, "owner", "admin", "finance", "auditor", "manager", "director", "procurement") {
		return forbidden(c)
	}

	canViewAll := middleware.HasAnyRole(user, "owner", "admin", "finance", "auditor", "director", "procurement")
	csvData, _, err := h.service.PurchaseRequestsCSV(c.UserContext(), user.OrganizationID, canViewAll, user.ID, ExportFilters{
		Status:       c.Query("status"),
		DepartmentID: c.Query("department_id"),
		FromDate:     c.Query("from_date"),
		ToDate:       c.Query("to_date"),
	})
	if err != nil {
		return writeError(c, err)
	}

	return sendCSV(c, "purchase-requests.csv", csvData)
}

func (h *Handler) approvalHistory(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}
	if !middleware.HasAnyRole(user, "owner", "admin", "finance", "auditor", "manager", "director", "procurement") {
		return forbidden(c)
	}

	csvData, _, err := h.service.ApprovalHistoryCSV(c.UserContext(), user.OrganizationID, user.ID, ExportFilters{
		RequestID: c.Query("request_id"),
		FromDate:  c.Query("from_date"),
		ToDate:    c.Query("to_date"),
	})
	if err != nil {
		return writeError(c, err)
	}

	return sendCSV(c, "approval-history.csv", csvData)
}

func (h *Handler) auditLogs(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}
	if !middleware.HasAnyRole(user, "owner", "admin", "finance", "auditor") {
		return forbidden(c)
	}

	csvData, _, err := h.service.AuditLogsCSV(c.UserContext(), user.OrganizationID, user.ID, ExportFilters{
		Action:     c.Query("action"),
		EntityType: c.Query("entity_type"),
		FromDate:   c.Query("from_date"),
		ToDate:     c.Query("to_date"),
	})
	if err != nil {
		return writeError(c, err)
	}

	return sendCSV(c, "audit-logs.csv", csvData)
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

func sendCSV(c *fiber.Ctx, filename, csvData string) error {
	c.Set(fiber.HeaderContentType, "text/csv; charset=utf-8")
	c.Set(fiber.HeaderContentDisposition, `attachment; filename="`+strings.ReplaceAll(filename, `"`, "")+`"`)
	return c.Status(fiber.StatusOK).SendString(csvData)
}
