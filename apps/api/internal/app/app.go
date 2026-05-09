package app

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"operra/api/internal/ai"
	"operra/api/internal/attachments"
	"operra/api/internal/audit"
	"operra/api/internal/auth"
	"operra/api/internal/departments"
	"operra/api/internal/exports"
	"operra/api/internal/platform/config"
	"operra/api/internal/platform/middleware"
	"operra/api/internal/platform/response"
	"operra/api/internal/requests"
	"operra/api/internal/storage"
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
	auditService := audit.NewService(a.db)
	auditHandler := audit.NewHandler(auditService)
	departmentHandler := departments.NewHandler(departments.NewService(a.db))
	userHandler := users.NewHandler(users.NewService(a.db))
	aiService := newAIService(a.cfg, a.db, auditService)
	workflowHandler := workflowapi.NewHandler(workflowapi.NewService(a.db, configworkflow.GenerateMermaid, configworkflow.ValidateConfig, auditService), aiService)
	requestService := requests.NewService(a.db, auditService)
	requestHandler := requests.NewHandler(requestService)
	exportHandler := exports.NewHandler(exports.NewService(a.db, auditService))
	attachmentHandler := attachments.NewHandler(attachments.NewService(a.db, buildStorage(a.cfg), auditService))

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
	requestsGroup := protectedAPI.Group("/purchase-requests")
	requestHandler.RegisterRoutes(requestsGroup)
	attachmentHandler.RegisterRoutes(requestsGroup)
	approvalsGroup := protectedAPI.Group("/approvals")
	approvalsGroup.Get("/pending", requestHandler.PendingApprovals)
	auditGroup := protectedAPI.Group("/audit-logs", middleware.RequireAnyRole("owner", "admin", "finance", "auditor"))
	auditHandler.RegisterRoutes(auditGroup)
	exportGroup := protectedAPI.Group("/exports")
	exportHandler.RegisterRoutes(exportGroup)
	departmentsGroup := protectedAPI.Group("/departments", middleware.RequireAnyRole("owner", "admin"))
	departmentHandler.RegisterRoutes(departmentsGroup)
	usersGroup := protectedAPI.Group("/users", middleware.RequireAnyRole("owner", "admin"))
	userHandler.RegisterRoutes(usersGroup)
}

func buildStorage(cfg config.Config) storage.Store {
	if cfg.StorageDriver == "s3" && cfg.S3Endpoint != "" && cfg.S3AccessKey != "" && cfg.S3SecretKey != "" {
		return storage.NewS3Store(cfg.S3Endpoint, cfg.S3Bucket, cfg.S3AccessKey, cfg.S3SecretKey, cfg.S3Region, cfg.S3ForcePathStyle)
	}
	return storage.NewLocalStore("")
}

func newAIService(cfg config.Config, db *sql.DB, auditService *audit.Service) *ai.Service {
	var provider ai.Provider
	if cfg.AIAPIKey != "" && cfg.AIBaseURL != "" {
		provider = ai.NewOpenAIProvider(cfg.AIBaseURL, cfg.AIAPIKey, cfg.AIModel)
	} else {
		provider = ai.NewRuleBasedProvider()
	}
	return ai.NewService(db, provider, auditService)
}
