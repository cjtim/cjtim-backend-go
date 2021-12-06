package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cjtim/cjtim-backend-go/configs"
	api "github.com/cjtim/cjtim-backend-go/internal/app/controllers"
	"github.com/cjtim/cjtim-backend-go/internal/app/middlewares"
	"github.com/cjtim/cjtim-backend-go/internal/app/repository"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	logger := middlewares.InitZap()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	app := startServer()
	setupCloseHandler(app)

	client := &repository.ClientImpl{}
	repository.Client = client
	err := client.Connect()
	if err != nil {
		zap.L().Error("Database start error", zap.Error(err))
		return 1
	}

	listen := fmt.Sprintf(":%d", configs.Config.Port)
	if err := app.Listen(listen); err != nil {
		repository.Client.Disconnect()
		zap.L().Error("fiber start error", zap.Error(err))
		return 1
	}
	return 0
}

func startServer() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: middlewares.ErrorHandling,
		BodyLimit:    100 * 1024 * 1024, // Limit file size to 100MB
	})
	app.Use(middlewares.Cors())
	app.Use(middlewares.RequestLog())
	api.Route(app) // setup router path
	return app
}

// setupCloseHandler - What to do when got ctrl+c SIGTERM
func setupCloseHandler(app *fiber.App) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	go func() {
		<-c
		zap.L().Info("Got SIGTERM, terminating program...")
		repository.Client.Disconnect()
		app.Server().Shutdown()
	}()
}
