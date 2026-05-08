package app

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"operra/api/internal/auth"
	"operra/api/internal/departments"
	"operra/api/internal/platform/config"
	"operra/api/internal/platform/middleware"
	"operra/api/internal/platform/response"
	"operra/api/internal/users"
	configworkflow "operra/api/internal/workflow"
	workflowapi "operra/api/internal/workflows"
)

type Logger interface {
	Printf(format string, args ...any)
}

type App struct {
	*fiber.App
	cfg config.Config
	db  *sql.DB
	log Logger
}

func New(cfg config.Config, db *sql.DB, log Logger) *App {
	f := fiber.New(fiber.Config{
		AppName:      "Operra API",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if fiberErr, ok := err.(*fiber.Error); ok {
				return response.Error(c, fiberErr.Code, response.APIError{
					Code:    http.StatusText(fiberErr.Code),
					Message: fiberErr.Message,
				})
			}

			return response.Error(c, fiber.StatusInternalServerError, response.APIError{
				Code:    "INTERNAL_ERROR",
				Message: "internal server error",
			})
		},
	})

	app := &App{App: f, cfg: cfg, db: db, log: log}
	app.routes()
	return app
}

func (a *App) routes() {
	authService := auth.NewService(a.db, a.cfg.JWTSecret)
	authHandler := auth.NewHandler(authService, a.cfg.JWTSecret)
	departmentHandler := departments.NewHandler(departments.NewService(a.db))
	userHandler := users.NewHandler(users.NewService(a.db))
	workflowHandler := workflowapi.NewHandler(workflowapi.NewService(a.db, configworkflow.GenerateMermaid, configworkflow.ValidateConfig))

	a.Get("/health", func(c *fiber.Ctx) error {
		return response.Success(c, fiber.Map{
			"status": "ok",
		})
	})

	a.Get("/ready", func(c *fiber.Ctx) error {
		if a.db == nil {
			return response.Error(c, fiber.StatusServiceUnavailable, response.APIError{
				Code:    "INTERNAL_ERROR",
				Message: "database not configured",
			})
		}

		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()

		if err := a.db.PingContext(ctx); err != nil {
			return response.Error(c, fiber.StatusServiceUnavailable, response.APIError{
				Code:    "INTERNAL_ERROR",
				Message: "database not ready",
			})
		}

		return response.Success(c, fiber.Map{
			"status": "ready",
		})
	})

	api := a.Group("/api/v1")
	authGroup := api.Group("/auth")
	authHandler.RegisterRoutes(authGroup)
	workflowGroup := api.Group("/workflows", middleware.RequireAuth(a.cfg.JWTSecret, authService), middleware.RequireAnyRole("owner", "admin"))
	workflowHandler.RegisterRoutes(workflowGroup)

	protectedAPI := api.Group("", middleware.RequireAuth(a.cfg.JWTSecret, authService))
	departmentsGroup := protectedAPI.Group("/departments", middleware.RequireAnyRole("owner", "admin"))
	departmentHandler.RegisterRoutes(departmentsGroup)
	usersGroup := protectedAPI.Group("/users", middleware.RequireAnyRole("owner", "admin"))
	userHandler.RegisterRoutes(usersGroup)
}
