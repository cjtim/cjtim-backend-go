package binance

import (
	"context"

	"github.com/cjtim/cjtim-backend-go/datasource/collections"
	"github.com/cjtim/cjtim-backend-go/models"
	"github.com/gofiber/fiber/v2"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.mongodb.org/mongo-driver/bson"
)

func Get(c *fiber.Ctx) error {
	data := collections.BinanceScheama{}
	user := c.Locals("user").(*linebot.UserProfileResponse)
	models := c.Locals("db").(*models.Models)
	err := models.FindOne("binance", &data, bson.M{"lineUid": user.UserID})
	if err != nil {
		return err
	}
	return c.JSON(data)
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
