package files

import (
	"github.com/cjtim/cjtim-backend-go/datasource/collections"
	"github.com/cjtim/cjtim-backend-go/models"
	"github.com/gofiber/fiber/v2"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.mongodb.org/mongo-driver/bson"
)

func List(c *fiber.Ctx) error {
	user := c.Locals("user").(*linebot.UserProfileResponse)
	models := c.Locals("db").(*models.Models)
	files := &[]collections.FileScheama{}
	if err := models.FindAll("files", files, bson.M{"lineUid": user.UserID}); err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"files": files,
	})
}

func Delete(c *fiber.Ctx) error {
	user := c.Locals("user").(*linebot.UserProfileResponse)
	body := &struct {
		Filename string `json:"fileName"`
	}{}
	if err := c.BodyParser(body); err != nil {
		return err
	}
	models := c.Locals("db").(*models.Models)
	files := &[]collections.FileScheama{}
	if err := models.Destroy("files", bson.M{"fileName": body.Filename, "lineUid": user.UserID}); err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"files": files,
	})
}
