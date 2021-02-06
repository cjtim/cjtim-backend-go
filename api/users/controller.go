package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/line/line-bot-sdk-go/linebot"
)

func Me(c *fiber.Ctx) error {
	profile := c.Locals("user").(*linebot.UserProfileResponse)
	return c.JSON(profile)
}
