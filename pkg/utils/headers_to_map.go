package utils

import "github.com/gofiber/fiber/v2"

func HeadersToMapStr(c *fiber.Ctx) map[string]string {
	headers := map[string]string{}
	c.Request().Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})
	return headers
}
