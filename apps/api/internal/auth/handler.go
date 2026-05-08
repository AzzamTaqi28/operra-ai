package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"operra/api/internal/platform/middleware"
	"operra/api/internal/platform/response"
)

type Handler struct {
	service *Service
	secret  string
}

func NewHandler(service *Service, secret string) *Handler {
	return &Handler{service: service, secret: secret}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Post("/register-organization", h.registerOrganization)
	router.Post("/login", h.login)
	router.Get("/me", middleware.RequireAuth(h.secret, h.service), h.me)
}

func (h *Handler) registerOrganization(c *fiber.Ctx) error {
	var req struct {
		OrganizationName string `json:"organization_name"`
		OrganizationSlug string `json:"organization_slug"`
		OwnerName        string `json:"owner_name"`
		OwnerEmail       string `json:"owner_email"`
		Password         string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, response.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "invalid request body",
		})
	}

	result, err := h.service.RegisterOrganization(c.UserContext(), RegisterOrganizationRequest{
		OrganizationName: req.OrganizationName,
		OrganizationSlug: req.OrganizationSlug,
		OwnerName:        req.OwnerName,
		OwnerEmail:       req.OwnerEmail,
		Password:         req.Password,
	})
	if err != nil {
		return writeServiceError(c, err)
	}

	return response.Success(c, result)
}

func (h *Handler) login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, response.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "invalid request body",
		})
	}

	result, err := h.service.Login(c.UserContext(), req)
	if err != nil {
		return writeServiceError(c, err)
	}

	return response.Success(c, result)
}

func (h *Handler) me(c *fiber.Ctx) error {
	user, ok := c.Locals("current_user").(*middleware.CurrentUser)
	if !ok || user == nil {
		return response.Error(c, fiber.StatusUnauthorized, response.APIError{
			Code:    "UNAUTHORIZED",
			Message: "invalid token",
		})
	}

	return response.Success(c, user)
}

func writeServiceError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, context.Canceled):
		return response.Error(c, fiber.StatusRequestTimeout, response.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "request canceled",
		})
	case err != nil && (err.Error() == "invalid credentials"):
		return response.Error(c, fiber.StatusUnauthorized, response.APIError{
			Code:    "UNAUTHORIZED",
			Message: "invalid email or password",
		})
	case err != nil && (err.Error() == "organization_name is required" || err.Error() == "organization_slug is required" || err.Error() == "owner_name is required" || err.Error() == "owner_email is required" || err.Error() == "password is required" || err.Error() == "email and password are required"):
		return response.Error(c, fiber.StatusBadRequest, response.APIError{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
	case err != nil && errors.Is(err, ErrConflict):
		return response.Error(c, fiber.StatusConflict, response.APIError{
			Code:    "CONFLICT",
			Message: "resource already exists",
		})
	case err != nil && strings.Contains(strings.ToLower(err.Error()), "conflict"):
		return response.Error(c, fiber.StatusConflict, response.APIError{
			Code:    "CONFLICT",
			Message: "resource already exists",
		})
	default:
		return response.Error(c, fiber.StatusInternalServerError, response.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "internal server error",
		})
	}
}
