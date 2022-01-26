package controller

import (
	"net/http"

	"github.com/cjtim/cjtim-backend-go/internal/app/controllers/binance"
	"github.com/cjtim/cjtim-backend-go/internal/app/controllers/files"
	line_controllers "github.com/cjtim/cjtim-backend-go/internal/app/controllers/line"
	"github.com/cjtim/cjtim-backend-go/internal/app/controllers/urls"
	"github.com/cjtim/cjtim-backend-go/internal/app/controllers/users"
	"github.com/cjtim/cjtim-backend-go/internal/app/middlewares"
	"github.com/cjtim/cjtim-backend-go/internal/app/repository"

	"github.com/gofiber/fiber/v2"
)

type RouteImpl struct {
	Files   FilesRoute
	Binance BinanceRoute
	Urls    UrlRoute
}

type FilesRoute struct {
	List_GET    string
	Upload_POST string
	Delete_POST string
}
type BinanceRoute struct {
	Get         string
	Wallet_GET  string
	Update_POST string
	CronJob_GET string
}

type UrlRoute struct {
	Add_POST    string
	List_GET    string
	Delete_POST string
}

var RoutePath = &RouteImpl{
	Files: FilesRoute{
		List_GET:    "/files/list",
		Upload_POST: "/files/upload",
		Delete_POST: "/files/delete",
	},
	Binance: BinanceRoute{
		Get:         "/binance/get",
		Wallet_GET:  "/binance/wallet",
		Update_POST: "/binance/update",
		CronJob_GET: "/binance/cronjob",
	},
	Urls: UrlRoute{
		Add_POST:    "/urls/add",
		List_GET:    "/urls/list",
		Delete_POST: "/urls/delete",
	},
}

// Route for all api request
func Route(r *fiber.App) {
	r.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"msg": "Hello, world"})
	})
	r.Get("/health", func(c *fiber.Ctx) error {
		if repository.Health() != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		return c.SendString("pong")
	})
	r.Post("/line/webhook", line_controllers.Webhook)
	r.Get("/line/weatherBroadcast", middlewares.InternalAuth, line_controllers.WeatherBroadcast)
	r.Get(RoutePath.Binance.CronJob_GET, middlewares.InternalAuth, binance.Cronjob)

	// Files
	r.Get(RoutePath.Files.List_GET, middlewares.LiffVerify, files.List)
	r.Post(RoutePath.Files.Upload_POST, middlewares.LiffVerify, files.Upload)
	r.Post(RoutePath.Files.Delete_POST, middlewares.LiffVerify, files.Delete)

	usersRouteSetup(r)

	// Binance
	r.Get(RoutePath.Binance.Get, middlewares.LiffVerify, binance.Get)
	r.Get(RoutePath.Binance.Wallet_GET, middlewares.LiffVerify, binance.GetWallet)
	r.Post(RoutePath.Binance.Update_POST, middlewares.LiffVerify, binance.UpdatePrice)

	// Urls
	r.Post(RoutePath.Urls.Add_POST, middlewares.LiffVerify, urls.Add)
	r.Get(RoutePath.Urls.List_GET, middlewares.LiffVerify, urls.List)
	r.Post(RoutePath.Urls.Delete_POST, middlewares.LiffVerify, urls.Delete)

}

func usersRouteSetup(r *fiber.App) {
	usersRoute := r.Group("/users", middlewares.LiffVerify)
	usersRoute.Get("/me", users.Me)
	usersRoute.Post("/update", users.Update)
}
