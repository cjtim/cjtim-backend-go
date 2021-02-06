package middlewares

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandling(c *fiber.Ctx, err error) error {
	log.Println(string(c.IP()), err)
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
