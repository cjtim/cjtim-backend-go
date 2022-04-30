package binance

import (
	"context"
	"net/http"

	"github.com/cjtim/cjtim-backend-go/internal/app/middlewares"
	"github.com/cjtim/cjtim-backend-go/internal/app/repository"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/binance"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/line_notify"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Get(c *fiber.Ctx) error {
	data := repository.BinanceScheama{}
	user := middlewares.GetUser(c)
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
	user := middlewares.GetUser(c)
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
	user := middlewares.GetUser(c)
	err = repository.BinanceRepo.FindOneAndReplace(context.TODO(), bson.M{"lineUid": user.UserID}, &data)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(200)
}

func Cronjob(c *fiber.Ctx) error {
	data := []repository.BinanceScheama{}
	err := repository.BinanceRepo.Find(&data, bson.M{})
	if err != nil {
		return err
	}

	_, errorUsers := line_notify.TriggerLineNotify(&data)

	if len(errorUsers) > 0 {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot send notify for some users.",
		})
	}

	return c.SendStatus(200)
}
