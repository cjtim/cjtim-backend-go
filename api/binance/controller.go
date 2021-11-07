package binance

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/cjtim/cjtim-backend-go/pkg/binance"
	"github.com/cjtim/cjtim-backend-go/pkg/utils"
	"github.com/cjtim/cjtim-backend-go/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func Get(c *fiber.Ctx) error {
	data := repository.BinanceScheama{}
	user := c.Locals("user").(*linebot.UserProfileResponse)
	collection := repository.DB.Collection("binance")
	result := collection.FindOne(context.TODO(), bson.M{"lineUid": user.UserID})
	userNotFound := result.Decode(&data)
	if userNotFound != nil {
		data = repository.BinanceScheama{
			LineUID:          user.UserID,
			LineNotifyToken:  "",
			BinanceApiKey:    "",
			BinanceSecretKey: "",
			Prices: map[string]interface{}{
				"BNB": 1,
			},
			LineNotifyTime: 5,
		}
		collection.InsertOne(context.TODO(), data)
		result := collection.FindOne(context.TODO(), bson.M{"lineUid": user.UserID})
		err := result.Decode(&data)
		if err != nil {
			return nil
		}
	}
	return c.JSON(data)
}

func GetWallet(c *fiber.Ctx) error {
	user := c.Locals("user").(*linebot.UserProfileResponse)
	userBinance := repository.BinanceScheama{}
	collection := repository.DB.Collection("binance")
	data := collection.FindOne(context.TODO(), bson.M{"lineUid": user.UserID})
	err := data.Decode(&userBinance)
	if err != nil {
		return nil
	}

	hasBinanceAPIKey := userBinance.BinanceApiKey != "" && userBinance.BinanceSecretKey != ""
	if hasBinanceAPIKey {
		binanceAccount, err := binance.GetBinanceAccount(userBinance.BinanceApiKey, userBinance.BinanceSecretKey)
		if err != nil {
			return err
		}
		return c.JSON(binanceAccount["balances"])
	}
	return c.JSON([]interface{}{})
}

func UpdatePrice(c *fiber.Ctx) error {
	data := repository.BinanceScheama{}
	err := c.BodyParser(&data)
	if err != nil {
		return err
	}
	user := c.Locals("user").(*linebot.UserProfileResponse)
	collection := repository.DB.Collection("binance")
	collection.FindOneAndReplace(context.TODO(), bson.M{"lineUid": user.UserID}, data)
	return c.SendStatus(200)
}

func Cronjob(c *fiber.Ctx) error {
	headers := utils.HeadersToMapStr(c)
	if headers["Authorization"] != os.Getenv("SECRET_PASSPHRASE") {
		return c.SendStatus(fiber.StatusForbidden)
	}
	collection := repository.DB.Collection("binance")
	data := &[]repository.BinanceScheama{}
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return err
	}
	err = cur.All(context.TODO(), data)
	if err != nil {
		return err
	}
	for _, user := range *data {
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
			_, _, _ = utils.Http(&utils.HttpReq{
				Method: http.MethodPost,
				URL:    os.Getenv("MICROSERVICE_BINANCE_LINE_NOTIFY_URL"),
				Headers: map[string]string{
					"Authorization": os.Getenv("SECRET_PASSPHRASE"),
				},
				Body: user,
			})
		}
	}
	return c.SendStatus(200)
}
