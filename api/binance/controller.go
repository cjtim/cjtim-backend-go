package binance

import (
	"context"
	"os"
	"time"

	"github.com/cjtim/cjtim-backend-go/datasource/collections"
	"github.com/cjtim/cjtim-backend-go/models"
	"github.com/cjtim/cjtim-backend-go/pkg/binance"
	"github.com/cjtim/cjtim-backend-go/pkg/utils"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.mongodb.org/mongo-driver/bson"
)

var restyClient = resty.New()

func Get(c *fiber.Ctx) error {
	data := collections.BinanceScheama{}
	user := c.Locals("user").(*linebot.UserProfileResponse)
	models := c.Locals("db").(*models.Models)
	collection := models.Client.Database("production").Collection("binance")
	result := collection.FindOne(context.TODO(), bson.M{"lineUid": user.UserID})
	userNotFound := result.Decode(&data)
	if userNotFound != nil {
		data = collections.BinanceScheama{
			LineUID:          user.UserID,
			LineNotifyToken:  "",
			BinanceApiKey:    "",
			BinanceSecretKey: "",
			Prices:           map[string]interface{}{},
			LineNotifyTime:   5,
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
	models := c.Locals("db").(*models.Models)
	userBinance := collections.BinanceScheama{}
	err := models.FindOne("binance", &userBinance, bson.M{"lineUid": user.UserID})
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
	data := collections.BinanceScheama{}
	err := c.BodyParser(&data)
	if err != nil {
		return err
	}
	user := c.Locals("user").(*linebot.UserProfileResponse)
	models := c.Locals("db").(*models.Models)
	collection := models.Client.Database("production").Collection("binance")
	collection.FindOneAndReplace(context.TODO(), bson.M{"lineUid": user.UserID}, data)
	return c.SendStatus(200)
}

func Cronjob(c *fiber.Ctx) error {
	headers := utils.HeadersToMapStr(c)
	if headers["Authorization"] != os.Getenv("SECRET_PASSPHRASE") {
		return c.SendStatus(fiber.StatusForbidden)
	}
	data := []collections.BinanceScheama{}
	models := c.Locals("db").(*models.Models)
	err := models.FindAll("binance", &data, nil)
	if err != nil {
		return err
	}
	for _, user := range data {
		userTime := user.LineNotifyTime % 60
		currentMinute := time.Now().Minute()
		needNotify := (currentMinute % int(userTime)) == 0
		if needNotify {
			restyClient.R().SetHeader("Authorization", os.Getenv("SECRET_PASSPHRASE")).SetBody(user).Post(
				os.Getenv("MICROSERVICE_BINANCE_LINE_NOTIFY_URL"),
			)
		}
	}
	return c.SendStatus(200)
}
