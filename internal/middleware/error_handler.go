package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Return JSON response
	return c.Status(code).JSON(fiber.Map{
		"status":  "error",
		"message": "Internal Server Error",
		"error":   err.Error(),
	})
}
