package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cjtim/cjtim-backend-go/api"
	"github.com/cjtim/cjtim-backend-go/datasource"
	"github.com/cjtim/cjtim-backend-go/middlewares"
	"github.com/cjtim/cjtim-backend-go/models"
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	app := startServer()
	if err := app.Listen(":8080"); err != nil {
		log.Panicln(err)
	}
}

func startServer() *fiber.App {
	var m *models.Models
	DBchannel := make(chan *mongo.Client)
	go datasource.MongoClient(DBchannel) // GoRoutine connectDB
	m = models.GetModels(<-DBchannel)
	app := fiber.New(fiber.Config{
		ErrorHandler: middlewares.ErrorHandling,
	})
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", m)
		return c.Next()
	})
	app.Use(middlewares.Cors())
	api.Route(app) // setup router path
	setupCloseHandler(app, m)
	return app
}

// setupCloseHandler - What to do when got ctrl+c SIGTERM
func setupCloseHandler(app *fiber.App, m *models.Models) {
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
