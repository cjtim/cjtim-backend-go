package line_controllers

import (
	"github.com/cjtim/cjtim-backend-go/config"
	"github.com/cjtim/cjtim-backend-go/pkg/airvisual"
	"github.com/cjtim/cjtim-backend-go/pkg/line"
	"github.com/cjtim/cjtim-backend-go/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// WeatherBroadcast - Broadcast Weather to line using CRON
func WeatherBroadcast(c *fiber.Ctx) error {
	headers := utils.HeadersToMapStr(c)
	if headers["Authorization"] != config.Config.SecretPassphrase {
		return c.SendStatus(fiber.StatusForbidden)
	}
	resp, err := airvisual.GetPhayaThaiCity()
	if err != nil {
		return err
	}
	msgs := []interface{}{line.WeatherFlexMessage(resp)}
	err = line.Broadcast(msgs)
	if err != nil {
		return nil
	}

	zap.L().Info("WeatherBroadcast")
	return c.SendStatus(fiber.StatusOK)
}
