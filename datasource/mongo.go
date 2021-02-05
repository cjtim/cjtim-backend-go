package datasource

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient for connect MongoDB
// GoRoutine
func MongoClient(datasourceChan chan *mongo.Client) error {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		return err
	}
	err = client.Connect(context.TODO())
	if err != nil {
		return err
	}
	if client.Ping(context.TODO(), nil) == nil {
		println("DB connected!")
	}
	datasourceChan <- client
	return nil
}
