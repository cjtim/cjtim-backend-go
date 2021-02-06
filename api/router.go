package api

import (
	"github.com/cjtim/cjtim-backend-go/api/files"
	"github.com/cjtim/cjtim-backend-go/api/line_webhook"
	"github.com/cjtim/cjtim-backend-go/api/users"
	"github.com/cjtim/cjtim-backend-go/middlewares"

	"github.com/gofiber/fiber/v2"
)

// Route for all api request
func Route(r *fiber.App) {
	r.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"msg": "Hello, world"})
	})
	r.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
	r.Post("/line/webhook", line_webhook.Webhook)
	// r.Post("/post", controllers.PostController)
	filesRouteSetup(r)
	userRouteSetup(r)
}

func filesRouteSetup(r *fiber.App) {
	fileRoute := r.Group("/files", middlewares.LiffVerify)
	fileRoute.Get("/list", files.List)
	fileRoute.Post("/upload", files.Upload)
	fileRoute.Post("/delete", nil)
}

func userRouteSetup(r *fiber.App) {
	usersRoute := r.Group("/users", middlewares.LiffVerify)
	usersRoute.Get("/me", users.Me)
}
