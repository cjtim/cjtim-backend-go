package binance

import (
	"context"
	"net/http"
	"time"

	"github.com/cjtim/cjtim-backend-go/configs"
	"github.com/cjtim/cjtim-backend-go/internal/app/repository"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/binance"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func Get(c *fiber.Ctx) error {
	data := repository.BinanceScheama{}
	user := c.Locals("user").(*linebot.UserProfileResponse)
	err := repository.BinanceRepo.FindOne(&data, bson.M{"lineUid": user.UserID})
	if err == mongo.ErrNoDocuments {
		data = repository.BinanceScheama{
			LineUID:          user.UserID,
			LineNotifyToken:  "",
			BinanceApiKey:    "",
			BinanceSecretKey: "",
			Prices: map[string]interface{}{
				"BNB": 1.00,
			},
			LineNotifyTime: 5,
		}
		_, err := repository.BinanceRepo.InsertOne(context.TODO(), &data)
		if err != nil {
			return err
		}
		return c.JSON(data)
	}
	if err != nil {
		return err
	}
	return c.JSON(data)
}

func GetWallet(c *fiber.Ctx) error {
	user := c.Locals("user").(*linebot.UserProfileResponse)
	userBinance := repository.BinanceScheama{}
	err := repository.BinanceRepo.FindOne(&userBinance, bson.M{"lineUid": user.UserID})
	if err != nil {
		return err
	}

	hasBinanceAPIKey := userBinance.BinanceApiKey != "" && userBinance.BinanceSecretKey != ""
	if hasBinanceAPIKey {
		binanceAccount, err := binance.GetBinanceAccount(userBinance.BinanceApiKey, userBinance.BinanceSecretKey)
		if err != nil {
			return err
		}
		return c.JSON(binanceAccount["balances"])
	}
	return c.SendStatus(http.StatusCreated)
}

func UpdatePrice(c *fiber.Ctx) error {
	data := repository.BinanceScheama{}
	err := c.BodyParser(&data)
	if err != nil {
		return err
	}
	user := c.Locals("user").(*linebot.UserProfileResponse)
	repository.BinanceRepo.FindOneAndReplace(context.TODO(), bson.M{"lineUid": user.UserID}, &data)
	return c.SendStatus(200)
}

func Cronjob(c *fiber.Ctx) error {
	headers := utils.HeadersToMapStr(c)
	_, found := headers["Authorization"]
	if !found || headers["Authorization"] != configs.Config.SecretPassphrase {
		return c.SendStatus(fiber.StatusForbidden)
	}
	data := []repository.BinanceScheama{}
	err := repository.BinanceRepo.Find(&data, bson.M{})
	if err != nil {
		return err
	}
	for _, user := range data {
		needNotify := false
		userTime := user.LineNotifyTime % 60
		currentMinute := time.Now().Minute()
		if userTime == 0 {
			needNotify = currentMinute == int(userTime)
		} else {
			needNotify = (currentMinute % int(userTime)) == 0
		}
		if needNotify {
			zap.L().Info("Trigger binance line notify", zap.String("lineuid", user.LineUID))
			resp, _, err := utils.Http(&utils.HttpReq{
				Method: http.MethodPost,
				URL:    configs.Config.LineNotifyURL,
				Headers: map[string]string{
					"Authorization": configs.Config.SecretPassphrase,
				},
				Body: user,
			})
			if err != nil || resp.StatusCode != http.StatusOK {
				zap.L().Error("Error trigger binance line notify", zap.String("lineuid", user.LineUID))
			}
		}
	}
	return c.SendStatus(200)
}
