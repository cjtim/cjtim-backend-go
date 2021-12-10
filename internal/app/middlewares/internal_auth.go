package middlewares

import (
	"github.com/cjtim/cjtim-backend-go/configs"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

var AuthorizationHeader = "Authorization"

func InternalAuth(c *fiber.Ctx) error {
	headers := utils.HeadersToMapStr(c)
	value, found := headers[AuthorizationHeader]
	if !found || value != configs.Config.SecretPassphrase {
		return c.SendStatus(fiber.StatusForbidden)
	}
	return nil
}
