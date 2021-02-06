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
	"github.com/cjtim/cjtim-backend-go/models"
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := startServer()
	if err := app.Listen(":8080"); err != nil {
		log.Panicln(err)
	}
}

func startServer() *fiber.App {
	m, err := models.GetModels(nil)
	if err != nil {
		log.Panic(err)
		return nil
	}
	app := fiber.New(fiber.Config{
		ErrorHandler: middlewares.ErrorHandling,
		BodyLimit:    100 * 1024 * 1024, // Limit file size to 4MB
	})
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", m)
		return c.Next()
	})
	app.Use(middlewares.Cors())
	api.Route(app) // setup router path
	setupCloseHandler(m)
	return app
}

// setupCloseHandler - What to do when got ctrl+c SIGTERM
func setupCloseHandler(m *models.Models) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Got SIGTERM, terminating program...")
		m.Client.Disconnect(context.TODO())
		fmt.Println("\r- MongoDB disconected!")
		os.Exit(0)
	}()
}
