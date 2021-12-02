package urls

import (
	"context"

	"github.com/cjtim/cjtim-backend-go/pkg/rebrandly"
	"github.com/cjtim/cjtim-backend-go/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.mongodb.org/mongo-driver/bson"
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
	insertID, err := repository.UrlRepo.InsertOne(context.TODO(), &resp)
	resp.ID = insertID
	if err != nil {
		return err
	}
	return c.JSON(resp)
}

func List(c *fiber.Ctx) error {
	urls := &[]repository.URLScheama{}
	user := c.Locals("user").(*linebot.UserProfileResponse)
	err := repository.UrlRepo.Find(&urls, bson.M{"lineUid": user.UserID})
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"urls": urls,
	})
}

func Delete(c *fiber.Ctx) error {
	user := c.Locals("user").(*linebot.UserProfileResponse)
	body := &URLDeleteBody{}
	err := c.BodyParser(body)
	if err != nil {
		return err
	}
	url := &repository.URLScheama{}
	err = repository.UrlRepo.FindOne(&url, bson.M{"lineUid": user.UserID, "shortUrl": body.ShortUrl})
	if err != nil {
		return err
	}
	err = rebrandly.Delete(url.RebrandlyID)
	if err != nil {
		return err
	}
	_, err = repository.UrlRepo.DeleteMany(context.TODO(), bson.M{
		"shortUrl": body.ShortUrl,
		"lineUid":  user.UserID,
	})
	if err != nil {
		return err
	}
	return nil
}
