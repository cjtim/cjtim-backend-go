package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ErrorHandling(c *fiber.Ctx, err error) error {
	zap.L().Error("Request errors", zap.String("ip", c.IP()), zap.Error(err))
	// Default 500 statuscode
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		// Override status code if fiber.Error type
		code = e.Code
	}
	// Set Content-Type: text/plain; charset=utf-8
	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	// Return statuscode with error message
	return c.Status(code).SendString(err.Error())
}
