package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cjtim/cjtim-backend-go/api"
	"github.com/cjtim/cjtim-backend-go/config"
	"github.com/cjtim/cjtim-backend-go/middlewares"
	"github.com/cjtim/cjtim-backend-go/repository"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var (
	app *fiber.App
)

func main() {
	logger := middlewares.InitZap()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	setupCloseHandler()

	_, err := repository.MongoClient()
	if err != nil {
		zap.L().Fatal("Database start error", zap.Error(err))
	}

	app = startServer()
	listen := fmt.Sprintf(":%d", config.Config.Port)
	if err := app.Listen(listen); err != nil {
		repository.DB.Client().Disconnect(context.TODO())
		zap.L().Info("MongoDB disconected!")
		zap.L().Fatal("fiber start error", zap.Error(err))
	}
	os.Exit(0)
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
func setupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	go func() {
		<-c
		zap.L().Info("Got SIGTERM, terminating program...")
		repository.Client.Disconnect(context.TODO())
		zap.L().Info("MongoDB disconected!")
		app.Server().Shutdown()
	}()
}
