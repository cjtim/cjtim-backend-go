package urls

import (
	"context"

	"github.com/cjtim/cjtim-backend-go/pkg/rebrandly"
	"github.com/cjtim/cjtim-backend-go/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URLAddBody struct {
	URL string `json:"url"`
}
type URLDeleteBody struct {
	ShortUrl string `json:"shortUrl"`
}

func Add(c *fiber.Ctx) error {
	user := c.Locals("user").(*linebot.UserProfileResponse)
	body := &URLAddBody{}
	err := c.BodyParser(body)
	if err != nil {
		return err
	}
	resp, err := rebrandly.Add(body.URL)
	if err != nil {
		return err
	}
	resp.LineUID = user.UserID
	insertID, err := repository.DB.Collection("urls").InsertOne(context.TODO(), resp)
	resp.ID = insertID.InsertedID.(primitive.ObjectID)
	if err != nil {
		return err
	}
	return c.JSON(resp)
}

func List(c *fiber.Ctx) error {
	urls := &[]repository.URLScheama{}
	user := c.Locals("user").(*linebot.UserProfileResponse)
	data, err := repository.DB.Collection("urls").Find(context.TODO(),bson.M{"lineUid": user.UserID})
	if err != nil {
		return err
	}
	err = data.All(context.TODO() ,urls)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"urls": urls,
	})
}

func Delete(c *fiber.Ctx) error {
	user := c.Locals("user").(*linebot.UserProfileResponse)
	collection := repository.DB.Collection("urls")
	body := &URLDeleteBody{}
	err := c.BodyParser(body)
	if err != nil {
		return err
	}
	url := &repository.URLScheama{}
	data := collection.FindOne(context.TODO(), bson.M{"lineUid": user.UserID, "shortUrl": body.ShortUrl})
	err = data.Decode(&url)
	if err != nil {
		return err
	}
	err = rebrandly.Delete(url.RebrandlyID)
	if err != nil {
		return err
	}
	_, err = collection.DeleteMany(context.TODO(), bson.M{
		"shortUrl": body.ShortUrl,
		"lineUid":  user.UserID,
	})
	if err != nil {
		return err
	}
	return nil
}
