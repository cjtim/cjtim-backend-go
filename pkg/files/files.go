package files

import (
	"mime"

	"github.com/cjtim/cjtim-backend-go/datasource/collections"
	"github.com/cjtim/cjtim-backend-go/models"
	"github.com/cjtim/cjtim-backend-go/pkg/gstorage"
	"github.com/cjtim/cjtim-backend-go/pkg/line"
	"github.com/gofiber/fiber/v2"
	"github.com/line/line-bot-sdk-go/linebot"
)

func Add(c *fiber.Ctx, e *linebot.Event) error {
	message := e.Message.(*linebot.FileMessage)
	fileByte, fileType, err := line.GetContent(message.ID)
	if err != nil {
		return err
	}
	ext, err := mime.ExtensionsByType(fileType)
	if ext == nil {
		ext = []string{""}
	}
	if err != nil {
		return err
	}
	Gclient, err := gstorage.GetClient()
	if err != nil {
		return err
	}
	objPath := ("users/" + e.Source.UserID + "/files/" + message.FileName + ext[0])
	url, err := Gclient.Upload(objPath, fileByte)
	if err != nil {
		return err
	}
	m := c.Locals("db").(*models.Models)
	data := &collections.FileScheama{
		FileName: message.FileName + ext[0],
		URL:      url,
		LineUID:  e.Source.UserID,
	}
	_, err = m.InsertOne("files", data)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"url": url,
	})
}
