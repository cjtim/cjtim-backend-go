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
		BodyLimit:    100 * 1024 * 1024, // Limit file size to 4MB
	})
	app.Use(func(c *fiber.Ctx) error {
		// 
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
		fmt.Println("\r- Got SIGTERM, terminating program...")
		repository.Client.Disconnect(context.TODO())
		fmt.Println("\r- MongoDB disconected!")
		os.Exit(0)
	}()
}
