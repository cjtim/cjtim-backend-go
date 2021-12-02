package users

import (
	"context"

	"github.com/cjtim/cjtim-backend-go/repository"
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
	body := repository.UserScheama{}

	err := c.BodyParser(&body)
	if err != nil {
		return err
	}
	userInDatabase, err := repository.UserRepo.CountDocuments(context.TODO(), profileFilter)
	if err != nil {
		return err
	}
	if userInDatabase < 1 {
		// new user
		body.LineUid = profile.UserID
		_, err := repository.UserRepo.InsertOne(context.TODO(), &body)
		if err != nil {
			return err
		}
	} else {
		_, err = repository.UserRepo.UpdateOne(context.TODO(), profileFilter, bson.M{"$set": body})
		if err != nil {
			return err
		}
	}
	return c.SendStatus(fiber.StatusOK)
}
