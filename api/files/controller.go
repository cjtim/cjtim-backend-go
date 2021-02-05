package files

import (
	"github.com/cjtim/cjtim-backend-go/datasource/collections"
	"github.com/cjtim/cjtim-backend-go/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func List(c *fiber.Ctx) error {
	models := c.Locals("db").(*models.Models)
	files := &[]collections.FileScheama{}
	if err := models.FindAll("files", files, bson.M{}); err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"files": files,
	})
}

func Delete(c *fiber.Ctx) error {
	body := &struct {
		Filename string `json:"fileName"`
	}{}
	if err := c.BodyParser(body); err != nil {
		return err
	}
	models := c.Locals("db").(*models.Models)
	files := &[]collections.FileScheama{}
	if err := models.Destroy("files", bson.M{"fileName": body.Filename}); err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"files": files,
	})
}
