package files

import (
	"io/ioutil"

	"github.com/cjtim/cjtim-backend-go/datasource/collections"
	"github.com/cjtim/cjtim-backend-go/models"
	"github.com/cjtim/cjtim-backend-go/pkg/files"
	"github.com/cjtim/cjtim-backend-go/pkg/gstorage"
	"github.com/gofiber/fiber/v2"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.mongodb.org/mongo-driver/bson"
)

func Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	data, err := file.Open()
	if err != nil {
		return err
	}
	bdata, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}
	user := c.Locals("user").(*linebot.UserProfileResponse)
	models := c.Locals("db").(*models.Models)
	files.Add(file.Filename, bdata, user.UserID, models)
	return nil
}

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
	err := c.BodyParser(body)
	if err != nil {
		return err
	}
	gClient, err := gstorage.GetClient()
	defer gClient.Client.Close()
	if err != nil {
		return err
	}
	err = gClient.Delete("users/" + user.UserID + "/files/" + body.Filename)
	if err != nil {
		return err
	}
	models := c.Locals("db").(*models.Models)
	err = models.Destroy("files", bson.M{"fileName": body.Filename, "lineUid": user.UserID})
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}
