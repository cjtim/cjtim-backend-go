package urls

import (
	"context"
	"net/http"

	"github.com/cjtim/cjtim-backend-go/internal/app/middlewares"
	"github.com/cjtim/cjtim-backend-go/internal/app/repository"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/rebrandly"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type URLAddBody struct {
	URL string `json:"url"`
}
type URLDeleteBody struct {
	ShortUrl string `json:"shortUrl"`
}

func Add(c *fiber.Ctx) error {
	user := middlewares.GetUser(c)
	body := URLAddBody{}
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}
	resp, err := rebrandly.Add(body.URL)
	if err != nil {
		return err
	}
	resp.LineUID = user.UserID
	insertID, err := repository.UrlRepo.InsertOne(context.TODO(), resp)
	resp.ID = insertID
	if err != nil {
		return err
	}
	return c.JSON(resp)
}

func List(c *fiber.Ctx) error {
	user := middlewares.GetUser(c)
	urls := []repository.URLScheama{}
	err := repository.UrlRepo.Find(&urls, bson.M{"lineUid": user.UserID})
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"urls": urls,
	})
}

func Delete(c *fiber.Ctx) error {
	user := middlewares.GetUser(c)
	body := URLDeleteBody{}
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}
	url := repository.URLScheama{}
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
	return c.SendStatus(http.StatusOK)
}
