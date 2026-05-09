package attachments

import (
	"database/sql"
	"errors"
	"io"
	"mime/multipart"
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
	router.Post("/:id/attachments", h.upload)
	router.Get("/:id/attachments/:attachment_id/download", h.download)
}

func (h *Handler) upload(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return badRequest(c, "file is required")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return writeError(c, err)
	}
	defer func() { _ = file.Close() }()

	data, err := io.ReadAll(file)
	if err != nil {
		return writeError(c, err)
	}

	attachment, err := h.service.Upload(c.UserContext(), user.OrganizationID, user.ID, user.Roles, c.Params("id"), fileHeader.Filename, fileHeader.Header.Get("Content-Type"), data)
	if err != nil {
		return writeError(c, err)
	}

	return response.Success(c, fiber.Map{
		"id":             attachment.ID,
		"file_name":      attachment.FileName,
		"file_size":      attachment.FileSize,
		"mime_type":      attachment.MIMEType,
		"created_at":     attachment.CreatedAt,
		"storage_key":    attachment.StorageKey,
		"storage_bucket": attachment.StorageBucket,
	})
}

func (h *Handler) download(c *fiber.Ctx) error {
	user, ok := middleware.CurrentUserFromCtx(c)
	if !ok {
		return unauthorized(c)
	}

	result, err := h.service.Download(c.UserContext(), user.OrganizationID, user.ID, user.Roles, c.Params("id"), c.Params("attachment_id"))
	if err != nil {
		return writeError(c, err)
	}

	c.Set(fiber.HeaderContentType, result.Attachment.MIMEType)
	c.Set(fiber.HeaderContentDisposition, `attachment; filename="`+strings.ReplaceAll(result.Attachment.FileName, `"`, "")+`"`)
	return c.Status(fiber.StatusOK).Send(result.Data)
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

func writeError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return response.Error(c, fiber.StatusNotFound, response.APIError{
			Code:    "NOT_FOUND",
			Message: "resource not found",
		})
	case errors.Is(err, ErrForbidden):
		return response.Error(c, fiber.StatusForbidden, response.APIError{
			Code:    "FORBIDDEN",
			Message: "forbidden",
		})
	default:
		return response.Error(c, fiber.StatusInternalServerError, response.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "internal server error",
		})
	}
}

func fileContentType(fileHeader *multipart.FileHeader) string {
	ct := fileHeader.Header.Get("Content-Type")
	if ct == "" {
		return "application/octet-stream"
	}
	return ct
}
