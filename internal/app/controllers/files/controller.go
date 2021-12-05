package files

import (
	"io/ioutil"

	"github.com/cjtim/cjtim-backend-go/internal/app/repository"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/files"
	"github.com/gofiber/fiber/v2"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.mongodb.org/mongo-driver/bson"
)

type Response struct {
	Files []repository.FileScheama `json:"files"`
}

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

	json, err := files.Client.Add(file.Filename, bdata, user.UserID)
	if err != nil {
		return err
	}
	return c.JSON(json)
}

func List(c *fiber.Ctx) error {
	user := c.Locals("user").(*linebot.UserProfileResponse)
	files := []repository.FileScheama{}
	err := repository.FileRepo.Find(&files, bson.M{"lineUid": user.UserID})
	if err != nil {
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
	err = files.Client.Delete(body.Filename, user.UserID)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}