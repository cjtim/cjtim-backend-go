package users

import "github.com/gofiber/fiber/v2"

func Me(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"users": "me",
	})
}
