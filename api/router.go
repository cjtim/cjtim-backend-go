package api

import (
	"github.com/cjtim/cjtim-backend-go/api/binance"
	"github.com/cjtim/cjtim-backend-go/api/files"
	line_controllers "github.com/cjtim/cjtim-backend-go/api/line"
	"github.com/cjtim/cjtim-backend-go/api/urls"
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
	r.Post("/line/webhook", line_controllers.Webhook)
	r.Get("/line/weatherBroadcast", line_controllers.WeatherBroadcast)
	r.Get("/binance/cronjob", binance.Cronjob)
	// r.Post("/post", controllers.PostController)
	filesRouteSetup(r)
	usersRouteSetup(r)
	urlsRouteSetup(r)
	binanceRouteSetup(r)
}

func filesRouteSetup(r *fiber.App) {
	fileRoute := r.Group("/files", middlewares.LiffVerify)
	fileRoute.Get("/list", files.List)
	fileRoute.Post("/upload", files.Upload)
	fileRoute.Post("/delete", files.Delete)
}

func usersRouteSetup(r *fiber.App) {
	usersRoute := r.Group("/users", middlewares.LiffVerify)
	usersRoute.Get("/me", users.Me)
}

func urlsRouteSetup(r *fiber.App) {
	urlsRoute := r.Group("/urls", middlewares.LiffVerify)
	urlsRoute.Post("/add", urls.Add)
	urlsRoute.Get("/list", urls.List)
	urlsRoute.Post("/delete", urls.Delete)
}

func binanceRouteSetup(r *fiber.App) {
	binanceRoute := r.Group("/binance", middlewares.LiffVerify)
	binanceRoute.Get("/get", binance.Get)
	binanceRoute.Post("/update", binance.UpdatePrice)
}
