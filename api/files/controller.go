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
	models.FindAll("files", files, bson.M{})
	return c.JSON(fiber.Map{
		"files": files,
	})
}
