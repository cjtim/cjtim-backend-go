package datasource

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient for connect MongoDB
// GoRoutine
func MongoClient(datasourceChan chan *mongo.Client) {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	if client.Ping(context.TODO(), nil) == nil {
		println("DB connected!")
	}
	datasourceChan <- client
}
