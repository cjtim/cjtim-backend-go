package api

import (
	"github.com/cjtim/cjtim-backend-go/internal/app/controllers/binance"
	"github.com/cjtim/cjtim-backend-go/internal/app/controllers/files"
	line_controllers "github.com/cjtim/cjtim-backend-go/internal/app/controllers/line"
	"github.com/cjtim/cjtim-backend-go/internal/app/controllers/urls"
	"github.com/cjtim/cjtim-backend-go/internal/app/controllers/users"
	"github.com/cjtim/cjtim-backend-go/internal/app/middlewares"

	"github.com/gofiber/fiber/v2"
)

// Route for all api request
func Route(r *fiber.App) {
	r.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"msg": "Hello, world"})
	})
	r.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
	r.Post("/line/webhook", line_controllers.Webhook)
	r.Get("/line/weatherBroadcast", line_controllers.WeatherBroadcast)

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
	usersRoute.Post("/update", users.Update)
}

func urlsRouteSetup(r *fiber.App) {
	urlsRoute := r.Group("/urls", middlewares.LiffVerify)
	urlsRoute.Post("/add", urls.Add)
	urlsRoute.Get("/list", urls.List)
	urlsRoute.Post("/delete", urls.Delete)
}

func binanceRouteSetup(r *fiber.App) {
	binanceRoute := r.Group("/binance")
	binanceRoute.Get("/get", middlewares.LiffVerify, binance.Get)
	binanceRoute.Get("/wallet", middlewares.LiffVerify, binance.GetWallet)
	binanceRoute.Post("/update", middlewares.LiffVerify, binance.UpdatePrice)
	binanceRoute.Get("/cronjob", binance.Cronjob)
}