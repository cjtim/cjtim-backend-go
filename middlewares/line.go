package middlewares

import (
	"strings"

	"github.com/cjtim/cjtim-backend-go/pkg/line"
	"github.com/cjtim/cjtim-backend-go/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func LiffVerify(c *fiber.Ctx) error {
	headers := utils.HeadersToMapStr(c)
	split_bearer := strings.Split(headers["Authorization"], " ")
	if len(split_bearer) <= 1 {
		return c.SendStatus(fiber.StatusForbidden)
	}
	token := string(split_bearer[1])
	err := line.LineIsTokenValid(token)
	if err != nil {
		return c.SendStatus(fiber.StatusForbidden)
	}
	profile, err := line.LineGetProfile(token)
	if err != nil {
		return err
	}
	c.Locals("user", profile)
	return c.Next()
}
