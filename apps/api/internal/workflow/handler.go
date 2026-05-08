package workflow

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"

	"operra/api/internal/platform/response"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Post("/validate", h.validate)
}

func (h *Handler) validate(c *fiber.Ctx) error {
	var req struct {
		Type       string          `json:"type"`
		ConfigJSON json.RawMessage `json:"config_json"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, response.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "invalid request body",
		})
	}

	var cfg Config
	if len(req.ConfigJSON) > 0 {
		if err := json.Unmarshal(req.ConfigJSON, &cfg); err != nil {
			return response.Error(c, fiber.StatusBadRequest, response.APIError{
				Code:    "VALIDATION_ERROR",
				Message: "config_json must be valid JSON",
			})
		}
	}

	if req.Type != "" {
		cfg.Type = req.Type
	}

	result := ValidateConfig(cfg)
	return response.Success(c, fiber.Map{
		"validation":      result,
		"mermaid_diagram": GenerateMermaid(cfg),
	})
}
