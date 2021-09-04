package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cjtim/cjtim-backend-go/api"
	"github.com/cjtim/cjtim-backend-go/middlewares"
	"github.com/cjtim/cjtim-backend-go/repository"
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	client, err := repository.MongoClient()
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}
	repository.Client = client
	repository.DB = client.Database(os.Getenv("MONGO_DB"))
	app := startServer()
	if err := app.Listen(":8080"); err != nil {
		log.Panicln(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func startServer() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: middlewares.ErrorHandling,
		BodyLimit:    100 * 1024 * 1024, // Limit file size to 100MB
	})
	app.Use(func (c *fiber.Ctx) error {
		ips := c.IPs()
		isProxy := len(ips) > 0
		zap.L().Info("X-REAL-IP", zap.Strings("ips", ips))
		if (isProxy) {
			zap.L().Info("Request", 
				zap.String("ip", ips[len(ips)-1]),
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
			)
		} else {
			zap.L().Info("Request", 
				zap.String("ip", c.IP()),
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
			)
		}
		return c.Next()
	}) 
	app.Use(middlewares.Cors())
	api.Route(app) // setup router path
	setupCloseHandler()
	return app
}

// setupCloseHandler - What to do when got ctrl+c SIGTERM
func setupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		zap.L().Info("Got SIGTERM, terminating program...")
		repository.Client.Disconnect(context.TODO())
		zap.L().Info("MongoDB disconected!")
		os.Exit(0)
	}()
}
