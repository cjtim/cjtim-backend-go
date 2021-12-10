package line_controllers

import (
	"github.com/cjtim/cjtim-backend-go/internal/pkg/airvisual"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/line"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// WeatherBroadcast - Broadcast Weather to line using CRON
func WeatherBroadcast(c *fiber.Ctx) error {
	resp, err := airvisual.GetPhayaThaiCity()
	if err != nil {
		zap.L().Error("WeatherBroadcast error get airvisual", zap.Error(err))
		return err
	}
	msgs := line.WeatherFlexMessage(resp)
	err = line.Broadcast(msgs)
	if err != nil {
		zap.L().Error("WeatherBroadcast error line.Broadcast", zap.Error(err))
		return nil
	}
	zap.L().Info("WeatherBroadcast")
	return c.SendStatus(fiber.StatusOK)
}
