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

func CurrentUserFromCtx(c *fiber.Ctx) (*CurrentUser, bool) {
	user, ok := c.Locals("current_user").(*CurrentUser)
	return user, ok && user != nil
}

func HasAnyRole(user *CurrentUser, roles ...string) bool {
	if user == nil {
		return false
	}

	roleSet := make(map[string]struct{}, len(user.Roles))
	for _, role := range user.Roles {
		roleSet[role] = struct{}{}
	}

	for _, role := range roles {
		if _, ok := roleSet[role]; ok {
			return true
		}
	}

	return false
}

func RequireAnyRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := CurrentUserFromCtx(c)
		if !ok {
			return response.Error(c, fiber.StatusUnauthorized, response.APIError{
				Code:    "UNAUTHORIZED",
				Message: "unauthorized",
			})
		}

		if !HasAnyRole(user, roles...) {
			return response.Error(c, fiber.StatusForbidden, response.APIError{
				Code:    "FORBIDDEN",
				Message: "forbidden",
			})
		}

		return c.Next()
	}
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
