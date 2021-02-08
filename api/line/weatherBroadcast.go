package line_controllers

import (
	"os"

	"github.com/cjtim/cjtim-backend-go/pkg/airvisual"
	"github.com/cjtim/cjtim-backend-go/pkg/line"
	"github.com/cjtim/cjtim-backend-go/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
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
	_, err = line.Broadcast([]linebot.SendingMessage{line.WeatherFlexMessage(resp)})
	if err != nil {
		return nil
	}
	return c.SendStatus(fiber.StatusOK)
}
