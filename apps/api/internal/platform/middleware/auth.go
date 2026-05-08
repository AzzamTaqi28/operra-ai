package middleware

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"operra/api/internal/platform/response"
	"operra/api/internal/platform/security"
)

type CurrentUser struct {
	ID             string   `json:"id"`
	OrganizationID string   `json:"organization_id"`
	Name           string   `json:"name"`
	Email          string   `json:"email"`
	Roles          []string `json:"roles"`
	Status         string   `json:"status"`
}

type authRepository interface {
	UserByIDAndOrganization(userID, organizationID string) (*CurrentUser, error)
}

type AuthService interface {
	CurrentUser(userID, organizationID string) (*CurrentUser, error)
}

func RequireAuth(secret string, svc AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return response.Error(c, fiber.StatusUnauthorized, response.APIError{
				Code:    "UNAUTHORIZED",
				Message: "missing bearer token",
			})
		}

		claims, err := security.ParseToken(secret, strings.TrimPrefix(authHeader, "Bearer "))
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, response.APIError{
				Code:    "UNAUTHORIZED",
				Message: "invalid token",
			})
		}

		user, err := svc.CurrentUser(claims.UserID, claims.OrganizationID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return response.Error(c, fiber.StatusUnauthorized, response.APIError{
					Code:    "UNAUTHORIZED",
					Message: "invalid token",
				})
			}

			return response.Error(c, fiber.StatusInternalServerError, response.APIError{
				Code:    "INTERNAL_ERROR",
				Message: "internal server error",
			})
		}

		c.Locals("current_user", user)
		return c.Next()
	}
}
