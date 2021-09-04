package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cjtim/cjtim-backend-go/api"
	"github.com/cjtim/cjtim-backend-go/middlewares"
	"github.com/cjtim/cjtim-backend-go/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
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
	app.Use(logger.New(logger.Config{
		Next: func(c *fiber.Ctx) bool {
			ips := c.IPs()
			isProxy := len(ips) > 0
			log.Default().Printf("IPs: %s\n", ips)
			if (isProxy) {
				log.Default().Printf("%s - %s %s", ips[len(ips)-1], c.Method(), c.Path())
			}
			return isProxy
		},
		Output: os.Stdout,
		Format:     "[${time}] ${ip} ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006 | 15:04:05",
		TimeZone:   "Asia/Bangkok",
	}))
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
		fmt.Println("\r- Got SIGTERM, terminating program...")
		repository.Client.Disconnect(context.TODO())
		fmt.Println("\r- MongoDB disconected!")
		os.Exit(0)
	}()
}
