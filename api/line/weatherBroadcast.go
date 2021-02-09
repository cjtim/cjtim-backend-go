package line_controllers

import (
	"os"

	"github.com/cjtim/cjtim-backend-go/pkg/airvisual"
	"github.com/cjtim/cjtim-backend-go/pkg/line"
	"github.com/cjtim/cjtim-backend-go/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var _ = godotenv.Load()

// WeatherBroadcast - Broadcast Weather to line using CRON
func WeatherBroadcast(c *fiber.Ctx) error {
	headers := utils.HeadersToMapStr(c)
	if headers["Authorization"] != os.Getenv("SECRET_PASSPHRASE") {
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
	return c.SendStatus(fiber.StatusOK)
}
