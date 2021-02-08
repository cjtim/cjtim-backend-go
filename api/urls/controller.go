package urls

import (
	"github.com/cjtim/cjtim-backend-go/datasource/collections"
	"github.com/cjtim/cjtim-backend-go/models"
	"github.com/cjtim/cjtim-backend-go/pkg/rebrandly"
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
	m := c.Locals("db").(*models.Models)
	resp.LineUID = user.UserID
	insertID, err := m.InsertOne("urls", resp)
	resp.ID = insertID
	if err != nil {
		return err
	}
	return c.JSON(resp)
}

func List(c *fiber.Ctx) error {
	m := c.Locals("db").(*models.Models)
	urls := &[]collections.URLScheama{}
	user := c.Locals("user").(*linebot.UserProfileResponse)
	err := m.FindAll("urls", urls, bson.M{"lineUid": user.UserID})
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"urls": urls,
	})
}

func Delete(c *fiber.Ctx) error {
	m := c.Locals("db").(*models.Models)
	user := c.Locals("user").(*linebot.UserProfileResponse)
	body := &URLDeleteBody{}
	err := c.BodyParser(body)
	if err != nil {
		return err
	}
	url := &collections.URLScheama{}
	err = m.FindOne("urls", url, bson.M{"lineUid": user.UserID, "shortUrl": body.ShortUrl})
	if err != nil {
		return err
	}
	err = rebrandly.Delete(url.RebrandlyID)
	if err != nil {
		return err
	}
	m.Destroy("urls", bson.M{
		"shortUrl": body.ShortUrl,
		"lineUid":  user.UserID,
	})
	return nil
}
