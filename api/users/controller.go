package users

import (
	"context"

	"github.com/cjtim/cjtim-backend-go/datasource/collections"
	"github.com/cjtim/cjtim-backend-go/models"
	"github.com/gofiber/fiber/v2"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.mongodb.org/mongo-driver/bson"
)

func Me(c *fiber.Ctx) error {
	profile := c.Locals("user").(*linebot.UserProfileResponse)
	return c.JSON(profile)
}

func Update(c *fiber.Ctx) error {
	profile := c.Locals("user").(*linebot.UserProfileResponse)
	profileFilter := bson.M{"lineUid": profile.UserID}
	db := c.Locals("db").(*models.Models)
	collection := db.Client.Database("production").Collection("users")
	body := collections.UserScheama{}

	err := c.BodyParser(&body)
	if err != nil {
		return err
	}
	userInDatabase, err := collection.CountDocuments(context.TODO(), profileFilter)
	if err != nil {
		return err
	}
	if userInDatabase < 1 {
		// new user
		body.LineUid = profile.UserID
		_, err := db.InsertOne("users", body)
		if err != nil {
			return err
		}
	} else {
		_, err = db.Update("users", body, profileFilter)
		if err != nil {
			return err
		}
	}
	return c.SendStatus(fiber.StatusOK)
}
