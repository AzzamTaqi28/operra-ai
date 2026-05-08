package response

import "github.com/gofiber/fiber/v2"

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func Success(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": data,
	})
}

func SuccessList(c *fiber.Ctx, data any, pagination any) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":       data,
		"pagination": pagination,
	})
}

func Error(c *fiber.Ctx, status int, err APIError) error {
	return c.Status(status).JSON(fiber.Map{
		"error": err,
	})
}
